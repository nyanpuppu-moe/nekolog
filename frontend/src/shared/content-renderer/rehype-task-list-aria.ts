import type { Root as HRoot } from "hast";
import { visit } from "unist-util-visit";

export function rehypeTaskListAria() {
  return (tree: HRoot) => {
    let count = 0;

    visit(tree, "element", (node) => {
      // Find gfm task list item
      if (
        node.tagName === "li" &&
        node.properties?.className?.toString()?.includes("task-list-item")
      ) {
        // Find input checkbox
        const inputIndex = node.children.findIndex(
          (c: any) =>
            c.tagName === "input" && c.properties?.type === "checkbox",
        );

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
