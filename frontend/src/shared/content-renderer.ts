import { unified } from "unified";
import remarkParse from "remark-parse";
import remarkRehype from "remark-rehype";
import remarkMath from "remark-math";
import remarkGfm from "remark-gfm";
import rehypeKatex from "rehype-katex";
import rehypeStringify from "rehype-stringify";
import rehypeHighlight from "rehype-highlight";
import remarkFrontmatter from "remark-frontmatter";
import { rehypeAccessibleEmojis } from "rehype-accessible-emojis";
import type { Root as MdRoot } from "mdast";
import type { Root as HRoot } from "hast";
import type { VFile } from "vfile";
import { visit } from "unist-util-visit";
import { sanitize, clearWindow } from "isomorphic-dompurify";

export function remarkLatexToMath() {
  return (tree: MdRoot) => {
    visit(tree, "code", (node) => {
      if (node.lang === "latex") {
        node.lang = "math";
      }
    });
  };
}

export function rehypeTrimTaskSpace() {
  return (tree: HRoot) => {
    visit(tree, "element", (node) => {
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

export function rehypeTaskListAria() {
  return (tree: HRoot) => {
    let count = 0;

    visit(tree, "element", (node) => {
      // Find gfm task list item
      console.log(node);
      if (
        node.tagName === "li" &&
        node.properties?.className?.toString()?.includes("task-list-item")
      ) {
        // Find input checkbox
        const inputIndex = node.children.findIndex(
          (c: any) =>
            c.tagName === "input" && c.properties?.type === "checkbox",
        );
        console.log(inputIndex);

        if (inputIndex !== -1) {
          const inputNode = node.children[inputIndex] as any;
          const uniqueId = `task-item-${++count}`;

          // Collect tailing texts after checkbox
          const remaining = node.children.slice(inputIndex + 1);

          // Find block element, not inline children
          const firstBlockIdx = remaining.findIndex(
            (c: any) =>
              c.type === "element" &&
              !["span", "strong", "em", "code", "a"].includes(c.tagName),
          );

          // Split inline part and block part
          const inlinePart =
            firstBlockIdx === -1
              ? remaining
              : remaining.slice(0, firstBlockIdx);
          const blockPart =
            firstBlockIdx === -1 ? [] : remaining.slice(firstBlockIdx);

          inputNode.properties = inputNode.properties || {};

          if (inlinePart.length > 0) {
            // Wrap inline part with span
            const spanNode = {
              type: "element",
              tagName: "span",
              properties: { id: uniqueId },
              children: inlinePart,
            };
            inputNode.properties["aria-labelledby"] = uniqueId;

            node.children = [
              ...node.children.slice(0, inputIndex + 1),
              spanNode,
              ...blockPart, // Block elements are placed in outside of span
            ] as any;
          } else if (
            blockPart.length > 0 &&
            (blockPart[0] as any).tagName === "p"
          ) {
            // If block element is first child, add tag to block element
            (blockPart[0] as any).properties =
              (blockPart[0] as any).properties || {};
            (blockPart[0] as any).properties.id = uniqueId;
            inputNode.properties["aria-labelledby"] = uniqueId;
          }
        }
      }
    });
  };
}

export function remarkReadingTime() {
  return (_tree: MdRoot, file: VFile) => {
    // Get raw markdown content
    const rawMarkdown = String(file.value);

    // Remove frontmatter, codeblock, markdown symbols
    const cleanText = rawMarkdown
      .replace(/^---[\s\S]+?---/, "")
      .replace(/```[\s\S]*?```/g, "")
      .replace(/`([^`]+)`/g, "$1")
      .replace(/!?\[([^\]]+)\]\([^)]+\)/g, "$1")
      .replace(/[#*_\-~>+\d.]/g, "");

    // Get langth without spaces
    const finalLength = cleanText.replace(/\s+/g, "").length;
    const readingTime = Math.ceil(finalLength / 350);

    // Inject readingTime to frontmatter
    const data = file.data as Record<string, any>;
    data.frontmatter = data.frontmatter || {};
    data.frontmatter.readingTime = readingTime;
    data.frontmatter.wordCount = cleanText.length;
  };
}

export const processor = unified()
  .use(remarkParse) // Parse markdown
  .use(remarkFrontmatter) // Metadata
  .use(remarkReadingTime) // Get reading time
  .use(remarkGfm) // Tables, Autolinks, Tasklists, Strikethrough
  .use(remarkLatexToMath) // Change latex to math
  .use(remarkMath) // Parse math
  .use(remarkRehype) // Turn into HTML
  .use(rehypeTrimTaskSpace) // Remove front space from task list
  .use(rehypeTaskListAria) // Task list aira label
  .use(rehypeKatex) // Render latex
  .use(rehypeHighlight) // Highlight codeblock
  .use(rehypeAccessibleEmojis) // Add aria to emojis
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
