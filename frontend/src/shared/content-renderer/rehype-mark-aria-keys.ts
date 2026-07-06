import type { Root as HRoot, Properties } from "hast";
import { visit } from "unist-util-visit";

export function rehypeMarkAriaKeys() {
  return (tree: HRoot) => {
    visit(tree, "element", (node) => {
      if (node.tagName === "pre") {
        // In client runtime, we should inject aria-label
        node.properties ??= {} as Properties;
        node.properties["data-i18n-key"] = "code_snippet_scroll_area";
        node.properties["tabindex"] = "0";
      }
    });
  };
}
