import { Crepe } from "@milkdown/crepe";
import { Milkdown, MilkdownProvider, useEditor } from "@milkdown/react";
import { editorViewOptionsCtx } from "@milkdown/core";
import type React from "react";

const test = `
# Using Crepe Editor

Crepe is a powerful, feature-rich Markdown editor built on top of Milkdown. It provides a complete editing experience with a beautiful UI and extensive customization options.

## Why Choose Crepe?

***

* 🚀 **Ready to Use**: Works out of the box with sensible defaults
* 🎨 **Beautiful UI**: Modern design with multiple theme options
* 🔧 **Highly Customizable**: Extensive configuration options
* 📦 **Feature Complete**: Includes all essential Markdown editing features
* 🛠️ **Extensible**: Built on Milkdown's plugin system

\`\`\`latex
\\mu asdfq
\`\`\`

\`\`\`typescript
default export function name(asdf) {
  return this;
}
\`\`\`

| asdf | asdf | asdf |
| ---- | ---- | ---- |
| asdf | asdf | asdf |
| asdf | asdf | asdf |

* [ ] asdf
* [x] asdf
* [ ] asdf
* asdf
- asdf
2. asdf
3. asdf
4. asdf

[asdfasdf](https://qwreey.moe)
`;

export function CrepeEditor(): React.ReactElement {
  const { get } = useEditor((root) => {
    const crepe = new Crepe({
      root,
      defaultValue: test,
      features: {
        [Crepe.Feature.TopBar]: true,
        [Crepe.Feature.Toolbar]: false,
      },
    });

    // Turn off spellcheck by defualt
    crepe.editor.config((ctx) => {
      ctx.update(editorViewOptionsCtx, (prev) => ({
        ...prev,
        attributes: {
          ...prev.attributes,
          spellcheck: "false",
        },
      }));
    });

    return crepe;
  });

  return <Milkdown />;
}

export function MilkdownEditorWrapper(): React.ReactElement {
  return (
    <MilkdownProvider>
      <CrepeEditor />
    </MilkdownProvider>
  );
}
