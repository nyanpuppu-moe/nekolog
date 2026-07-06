import type { Root as HRoot } from "hast";
import { visit } from "unist-util-visit";

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
