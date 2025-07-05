import Routing from "./components/navigation/page";
// import Footer from "./footer/page";
import "./global.css";
import { Inter } from "next/font/google";
import {WebSocketProvider} from "@/app/context/WebSocketContext.jsx"

const inter = Inter({
  subsets: ["latin"],
  display: "swap",
  weight: ["100", "300", "400", "600", "700"],
});

export const metadata = {
  title: "Social network",
  description: "social network project",
};

export default function Layout({ children }) {
  return (
    <html lang="en" className={inter.className}>
      <head>
        <link
          href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined"
          rel="stylesheet"
        />
        <link rel="icon" href="/favicon.ico" />
      </head>
      <body>
        <WebSocketProvider>
          <Routing />
          <div>{children}</div>
          {/* <Footer/> */}
        </WebSocketProvider>
      </body>
    </html>
  );
}
