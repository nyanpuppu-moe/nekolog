import type { Root as MdRoot } from "mdast";
import { visit } from "unist-util-visit";

export function remarkLatexToMath() {
  return (tree: MdRoot) => {
    visit(tree, "code", (node) => {
      if (node.lang === "latex") {
        node.lang = "math";
      }
    });
  };
}
