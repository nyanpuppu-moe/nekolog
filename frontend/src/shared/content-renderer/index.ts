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
import { sanitize, clearWindow } from "isomorphic-dompurify";

import { remarkLatexToMath } from "./remark-latex-to-math";
import { rehypeTrimTaskSpace } from "./rehype-trim-task-space";
import { remarkReadingTime } from "./remark-reading-time";
import { rehypeTaskListAria } from "./rehype-task-list-aria";
import { rehypeMarkAriaKeys } from "./rehype-mark-aria-keys";

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
  .use(rehypeMarkAriaKeys) // Add pre aria locale keys
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
