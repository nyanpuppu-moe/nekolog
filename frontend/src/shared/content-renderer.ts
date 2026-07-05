import { unified } from "unified";
import remarkParse from "remark-parse";
import remarkRehype from "remark-rehype";
import remarkMath from "remark-math";
import remarkGfm from "remark-gfm";
import rehypeKatex from "rehype-katex";
import rehypeStringify from "rehype-stringify";
import rehypeHighlight from "rehype-highlight";
import type { Root } from "mdast";
import { visit } from "unist-util-visit";
import { sanitize, clearWindow } from "isomorphic-dompurify";

export function remarkLatexToMath() {
  return (tree: Root) => {
    visit(tree, "code", (node) => {
      if (node.lang === "latex") {
        node.lang = "math";
      }
    });
  };
}

export function rehypeTrimTaskSpace() {
  return (tree: any) => {
    visit(tree, "element", (node: any) => {
      // inspecat li
      if (node.tagName === "li") {
        const children = node.children;
        if (children && children.length > 1) {
          const first = children[0];
          const second = children[1];

          // first child is checkbox input
          if (
            first.type === "element" &&
            first.tagName === "input" &&
            first.properties?.type === "checkbox"
          ) {
            // if second child is text and starts with space, trim it
            if (second.type === "text") {
              if (second.value === " ") {
                children.splice(1, 1);
              } else if (second.value.startsWith(" ")) {
                second.value = second.value.slice(1);
              }
            }
          }
        }
      }
    });
  };
}

const processor = unified()
  .use(remarkParse) // Parse markdown
  .use(remarkGfm) // Tables, Autolinks, Tasklists, Strikethrough
  .use(remarkLatexToMath) // Change latex to math
  .use(remarkMath)
  .use(remarkRehype) // Turn into HTML
  .use(rehypeTrimTaskSpace) // Remove front space from task list
  .use(rehypeKatex)
  .use(rehypeHighlight)
  .use(rehypeStringify);

export async function compileAndSanitizeMarkdown(
  markdownRaw: string,
): Promise<string> {
  const rawHtml = await processor.process(markdownRaw); // Turn into string HTML

  // Sanitize content to prevent XSS
  const securedHtml = sanitize(String(rawHtml), {
    USE_PROFILES: { html: true, svg: true, mathMl: true },
    ADD_ATTR: ["data-project-tree", "contenteditable"],
  });
  clearWindow();

  return securedHtml;
}
