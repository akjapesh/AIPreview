export default function Layout({ children }) {
  return (
    <div
      style={{
        height: "100vh",
        width: "100vw",
        overflow: "hidden",
      }}
    >
      {children}
    </div>
  );
}
