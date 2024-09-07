import { useInput, useFieldValue } from "react-admin";
import MDEditor from "@uiw/react-md-editor";

export const MarkdownInput = ({ source, label }) => {
  const { id, field } = useInput({ source });
  return (
    <div id={id} data-color-mode="light" style={{ width: "100%" }}>
      <div className="wmde-markdown-var"> </div>
      <MDEditor {...field} height={500} />
    </div>
  );
};

export const MarkdownField = ({ source }) => {
  const value = useFieldValue({ source });
  if (!value) {
    return <span>-</span>;
  }

  return (
    <div data-color-mode="light" style={{ width: "100%" }}>
      <MDEditor.Markdown source={value} />
    </div>
  );
};
