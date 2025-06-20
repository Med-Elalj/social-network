import Routing from "./components/navigation/page";
// import Footer from "./footer/page";
import "./global.css";

export const metadata = {
  title: "Social network",
  description: "social network project",
};

export default function Layout({ children }) {
  return (
    <html lang="en">
      <body>
        <Routing />
        <div>{children}</div>
        {/* <Footer/> */}
      </body>
    </html>
  );
}
