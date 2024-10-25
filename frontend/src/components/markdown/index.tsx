import { useInput, useFieldValue } from "react-admin";
import MDEditor from "@uiw/react-md-editor";
import PropTypes from "prop-types";

export const MarkdownInput = ({ source }) => {
  const { id, field } = useInput({ source });
  return (
    <div id={id} data-color-mode="light" style={{ width: "100%" }}>
      <div className="wmde-markdown-var"> </div>
      <MDEditor {...field} height={500} />
    </div>
  );
};

MarkdownInput.propTypes = {
  source: PropTypes.string,
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

MarkdownField.propTypes = {
  source: PropTypes.string,
};
