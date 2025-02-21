import React from "react";
import { useStyletron } from "baseui";

export function Preview({ code }) {
  const [css] = useStyletron();
  const [previewUrl, setPreviewUrl] = useState("");

  useEffect(() => {
    // This is a simplified example. In a real app, you'd use a
    // sandboxed iframe solution like Sandpack or a custom preview service
    if (code) {
      // Placeholder for preview generation logic
      setPreviewUrl(`/api/preview?code=${encodeURIComponent(code)}`);
    }
  }, [code]);

  if (!previewUrl) {
    return <div>No preview available yet. Generate some code first!</div>;
  }

  return (
    <div className={css({ height: "100%" })}>
      <iframe
        src={previewUrl}
        className={css({
          width: "100%",
          height: "100%",
          border: "none",
          borderRadius: "4px",
        })}
        title="Component Preview"
        sandbox="allow-scripts"
      />
    </div>
  );
}
