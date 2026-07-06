import type { Root as MdRoot } from "mdast";
import type { VFile } from "vfile";

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
