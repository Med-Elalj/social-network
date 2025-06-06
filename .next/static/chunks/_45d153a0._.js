(globalThis.TURBOPACK = globalThis.TURBOPACK || []).push([typeof document === "object" ? document.currentScript : undefined, {

"[project]/src/app/navigation/nav.module.css [app-client] (css module)": ((__turbopack_context__) => {

var { g: global, __dirname } = __turbopack_context__;
{
__turbopack_context__.v({
  "active": "nav-module__ZdNNxq__active",
  "authLinks": "nav-module__ZdNNxq__authLinks",
  "bottomNav": "nav-module__ZdNNxq__bottomNav",
  "centerSection": "nav-module__ZdNNxq__centerSection",
  "dropdownMenu": "nav-module__ZdNNxq__dropdownMenu",
  "dropdownWrapper": "nav-module__ZdNNxq__dropdownWrapper",
  "iconDefault": "nav-module__ZdNNxq__iconDefault",
  "iconHover": "nav-module__ZdNNxq__iconHover",
  "iconWrapper": "nav-module__ZdNNxq__iconWrapper",
  "leftSection": "nav-module__ZdNNxq__leftSection",
  "linkWithIcon": "nav-module__ZdNNxq__linkWithIcon",
  "links": "nav-module__ZdNNxq__links",
  "logged": "nav-module__ZdNNxq__logged",
  "loginTitle": "nav-module__ZdNNxq__loginTitle",
  "nav": "nav-module__ZdNNxq__nav",
  "navLinks": "nav-module__ZdNNxq__navLinks",
  "notif": "nav-module__ZdNNxq__notif",
  "profile": "nav-module__ZdNNxq__profile",
  "rightSection": "nav-module__ZdNNxq__rightSection",
});
}}),
"[project]/utils/sendData.js [app-client] (ecmascript)": ((__turbopack_context__) => {
"use strict";

var { g: global, __dirname, k: __turbopack_refresh__, m: module } = __turbopack_context__;
{
__turbopack_context__.s({
    "SendData": (()=>SendData)
});
async function SendData(url, Data) {
    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(Data)
        });
        const responseBody = await response.json();
        return {
            status: response.status,
            body: responseBody
        };
    } catch (error) {
        return;
    }
}
_c = SendData;
var _c;
__turbopack_context__.k.register(_c, "SendData");
if (typeof globalThis.$RefreshHelpers$ === 'object' && globalThis.$RefreshHelpers !== null) {
    __turbopack_context__.k.registerExports(module, globalThis.$RefreshHelpers$);
}
}}),
"[project]/EndPoints/Logout.js [app-client] (ecmascript)": ((__turbopack_context__) => {
"use strict";

var { g: global, __dirname, k: __turbopack_refresh__, m: module } = __turbopack_context__;
{
__turbopack_context__.s({
    "Logout": (()=>Logout)
});
var __TURBOPACK__imported__module__$5b$project$5d2f$utils$2f$sendData$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__ = __turbopack_context__.i("[project]/utils/sendData.js [app-client] (ecmascript)");
;
async function Logout() {
    try {
        const { status, body } = await (0, __TURBOPACK__imported__module__$5b$project$5d2f$utils$2f$sendData$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["SendData"])('/api/v1/auth/logout', null);
        if (status === 200) {
            console.log("Logout successful:", body);
            // Clear localStorage/sessionStorage or cookies
            localStorage.removeItem("UserInfo");
            // Redirect to home or login
            window.location.href = "/auth/login"; // or "/" if you prefer
        } else {
            console.error("Logout failed with status", status);
        }
    } catch (err) {
        console.error("Logout error:", err);
    }
}
_c = Logout;
var _c;
__turbopack_context__.k.register(_c, "Logout");
if (typeof globalThis.$RefreshHelpers$ === 'object' && globalThis.$RefreshHelpers !== null) {
    __turbopack_context__.k.registerExports(module, globalThis.$RefreshHelpers$);
}
}}),
"[project]/src/app/navigation/page.jsx [app-client] (ecmascript)": ((__turbopack_context__) => {
"use strict";

var { g: global, __dirname, k: __turbopack_refresh__, m: module } = __turbopack_context__;
{
__turbopack_context__.s({
    "default": (()=>Routing)
});
var __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__ = __turbopack_context__.i("[project]/node_modules/next/dist/compiled/react/jsx-dev-runtime.js [app-client] (ecmascript)");
var __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$navigation$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__ = __turbopack_context__.i("[project]/node_modules/next/navigation.js [app-client] (ecmascript)");
var __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$index$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__ = __turbopack_context__.i("[project]/node_modules/next/dist/compiled/react/index.js [app-client] (ecmascript)");
var __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$client$2f$app$2d$dir$2f$link$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__ = __turbopack_context__.i("[project]/node_modules/next/dist/client/app-dir/link.js [app-client] (ecmascript)");
var __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$image$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__ = __turbopack_context__.i("[project]/node_modules/next/image.js [app-client] (ecmascript)");
var __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$app$2f$navigation$2f$nav$2e$module$2e$css__$5b$app$2d$client$5d$__$28$css__module$29$__ = __turbopack_context__.i("[project]/src/app/navigation/nav.module.css [app-client] (css module)");
var __TURBOPACK__imported__module__$5b$project$5d2f$utils$2f$sendData$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__ = __turbopack_context__.i("[project]/utils/sendData.js [app-client] (ecmascript)");
var __TURBOPACK__imported__module__$5b$project$5d2f$EndPoints$2f$Logout$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__ = __turbopack_context__.i("[project]/EndPoints/Logout.js [app-client] (ecmascript)");
;
var _s = __turbopack_context__.k.signature();
"use client";
;
;
;
;
;
;
;
function Routing() {
    _s();
    const [isOpen, setIsOpen] = (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$index$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["useState"])(false);
    const [isLoggedIn, setIsLoggedIn] = (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$index$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["useState"])(true);
    const pathname = (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$navigation$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["usePathname"])();
    //TODO: if logged forbedin to use login/logout and the oposet true
    (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$index$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["useEffect"])({
        "Routing.useEffect": ()=>{
            const fetchAuthStatus = {
                "Routing.useEffect.fetchAuthStatus": async ()=>{
                    try {
                        const { status, body } = await (0, __TURBOPACK__imported__module__$5b$project$5d2f$utils$2f$sendData$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["SendData"])('/api/v1/auth', null);
                        if (status === 200) {
                            setIsLoggedIn(true);
                            const { status: st, body: body2 } = await (0, __TURBOPACK__imported__module__$5b$project$5d2f$utils$2f$sendData$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["SendData"])('/api/v1/profile', {
                                id: body.id
                            });
                            if (st === 200) {
                                localStorage.setItem("UserInfo", JSON.stringify(body2.Userinfo));
                            } else {
                                console.log("err", st, body2);
                            }
                        }
                    } catch (err) {
                        console.error("Auth check failed", err);
                    }
                }
            }["Routing.useEffect.fetchAuthStatus"];
            fetchAuthStatus();
        }
    }["Routing.useEffect"], []);
    //TODO:fucntion to get data of user from local storage
    return /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])("div", {
        children: [
            /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])("div", {
                className: __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$app$2f$navigation$2f$nav$2e$module$2e$css__$5b$app$2d$client$5d$__$28$css__module$29$__["default"].nav,
                children: [
                    /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])("div", {
                        className: __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$app$2f$navigation$2f$nav$2e$module$2e$css__$5b$app$2d$client$5d$__$28$css__module$29$__["default"].leftSection,
                        children: /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])(__TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$client$2f$app$2d$dir$2f$link$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["default"], {
                            className: __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$app$2f$navigation$2f$nav$2e$module$2e$css__$5b$app$2d$client$5d$__$28$css__module$29$__["default"].loginTitle,
                            href: "/",
                            children: "Social Network"
                        }, void 0, false, {
                            fileName: "[project]/src/app/navigation/page.jsx",
                            lineNumber: 46,
                            columnNumber: 11
                        }, this)
                    }, void 0, false, {
                        fileName: "[project]/src/app/navigation/page.jsx",
                        lineNumber: 45,
                        columnNumber: 9
                    }, this),
                    /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])("div", {
                        className: __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$app$2f$navigation$2f$nav$2e$module$2e$css__$5b$app$2d$client$5d$__$28$css__module$29$__["default"].centerSection,
                        children: isLoggedIn && /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])(__TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["Fragment"], {
                            children: [
                                /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])(NavLink, {
                                    href: "/",
                                    icon: "home",
                                    pathname: pathname
                                }, void 0, false, {
                                    fileName: "[project]/src/app/navigation/page.jsx",
                                    lineNumber: 53,
                                    columnNumber: 15
                                }, this),
                                /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])(NavLink, {
                                    href: "/posts",
                                    icon: "posts",
                                    pathname: pathname
                                }, void 0, false, {
                                    fileName: "[project]/src/app/navigation/page.jsx",
                                    lineNumber: 54,
                                    columnNumber: 15
                                }, this),
                                /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])(NavLink, {
                                    href: "/groups",
                                    icon: "groupe",
                                    pathname: pathname
                                }, void 0, false, {
                                    fileName: "[project]/src/app/navigation/page.jsx",
                                    lineNumber: 55,
                                    columnNumber: 15
                                }, this),
                                /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])(NavLink, {
                                    href: "/chat",
                                    icon: "messages",
                                    pathname: pathname
                                }, void 0, false, {
                                    fileName: "[project]/src/app/navigation/page.jsx",
                                    lineNumber: 56,
                                    columnNumber: 15
                                }, this)
                            ]
                        }, void 0, true)
                    }, void 0, false, {
                        fileName: "[project]/src/app/navigation/page.jsx",
                        lineNumber: 50,
                        columnNumber: 9
                    }, this),
                    /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])("div", {
                        className: __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$app$2f$navigation$2f$nav$2e$module$2e$css__$5b$app$2d$client$5d$__$28$css__module$29$__["default"].rightSection,
                        children: isLoggedIn ? /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])(__TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["Fragment"], {
                            children: [
                                /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])("div", {
                                    className: __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$app$2f$navigation$2f$nav$2e$module$2e$css__$5b$app$2d$client$5d$__$28$css__module$29$__["default"].dropdownWrapper,
                                    children: /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])("div", {
                                        className: __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$app$2f$navigation$2f$nav$2e$module$2e$css__$5b$app$2d$client$5d$__$28$css__module$29$__["default"].notif,
                                        onClick: ()=>setIsOpen(true),
                                        onMouseLeave: ()=>setIsOpen(false),
                                        children: [
                                            /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])("span", {
                                                children: /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])(__TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$image$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["default"], {
                                                    src: "/notification.svg",
                                                    alt: "notification",
                                                    width: 25,
                                                    height: 25
                                                }, void 0, false, {
                                                    fileName: "[project]/src/app/navigation/page.jsx",
                                                    lineNumber: 72,
                                                    columnNumber: 21
                                                }, this)
                                            }, void 0, false, {
                                                fileName: "[project]/src/app/navigation/page.jsx",
                                                lineNumber: 71,
                                                columnNumber: 19
                                            }, this),
                                            isOpen && /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])("div", {
                                                className: __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$app$2f$navigation$2f$nav$2e$module$2e$css__$5b$app$2d$client$5d$__$28$css__module$29$__["default"].dropdownMenu,
                                                children: [
                                                    /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])(__TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$client$2f$app$2d$dir$2f$link$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["default"], {
                                                        href: `/`,
                                                        onClick: ()=>setIsOpen(false),
                                                        children: "test1"
                                                    }, void 0, false, {
                                                        fileName: "[project]/src/app/navigation/page.jsx",
                                                        lineNumber: 76,
                                                        columnNumber: 23
                                                    }, this),
                                                    /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])(__TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$client$2f$app$2d$dir$2f$link$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["default"], {
                                                        href: `/`,
                                                        onClick: ()=>setIsOpen(false),
                                                        children: "test2"
                                                    }, void 0, false, {
                                                        fileName: "[project]/src/app/navigation/page.jsx",
                                                        lineNumber: 77,
                                                        columnNumber: 23
                                                    }, this)
                                                ]
                                            }, void 0, true, {
                                                fileName: "[project]/src/app/navigation/page.jsx",
                                                lineNumber: 75,
                                                columnNumber: 21
                                            }, this)
                                        ]
                                    }, void 0, true, {
                                        fileName: "[project]/src/app/navigation/page.jsx",
                                        lineNumber: 66,
                                        columnNumber: 17
                                    }, this)
                                }, void 0, false, {
                                    fileName: "[project]/src/app/navigation/page.jsx",
                                    lineNumber: 65,
                                    columnNumber: 15
                                }, this),
                                /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])("div", {
                                    className: __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$app$2f$navigation$2f$nav$2e$module$2e$css__$5b$app$2d$client$5d$__$28$css__module$29$__["default"].dropdownWrapper,
                                    children: /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])("div", {
                                        className: __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$app$2f$navigation$2f$nav$2e$module$2e$css__$5b$app$2d$client$5d$__$28$css__module$29$__["default"].profile,
                                        onClick: ()=>setIsOpen(true),
                                        onMouseLeave: ()=>setIsOpen(false),
                                        children: [
                                            /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])("span", {
                                                className: __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$app$2f$navigation$2f$nav$2e$module$2e$css__$5b$app$2d$client$5d$__$28$css__module$29$__["default"].iconUser,
                                                children: /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])(__TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$image$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["default"], {
                                                    src: "/iconMale.png",
                                                    alt: "profile",
                                                    width: 40,
                                                    height: 40
                                                }, void 0, false, {
                                                    fileName: "[project]/src/app/navigation/page.jsx",
                                                    lineNumber: 89,
                                                    columnNumber: 21
                                                }, this)
                                            }, void 0, false, {
                                                fileName: "[project]/src/app/navigation/page.jsx",
                                                lineNumber: 88,
                                                columnNumber: 19
                                            }, this),
                                            /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])("span", {
                                                children: "nickname"
                                            }, void 0, false, {
                                                fileName: "[project]/src/app/navigation/page.jsx",
                                                lineNumber: 91,
                                                columnNumber: 19
                                            }, this),
                                            isOpen && /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])("div", {
                                                className: __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$app$2f$navigation$2f$nav$2e$module$2e$css__$5b$app$2d$client$5d$__$28$css__module$29$__["default"].dropdownMenu,
                                                children: [
                                                    /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])(__TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$client$2f$app$2d$dir$2f$link$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["default"], {
                                                        href: `/profile/nickname`,
                                                        onClick: ()=>setIsOpen(false),
                                                        children: "Profile"
                                                    }, void 0, false, {
                                                        fileName: "[project]/src/app/navigation/page.jsx",
                                                        lineNumber: 94,
                                                        columnNumber: 23
                                                    }, this),
                                                    /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])("button", {
                                                        onClick: async ()=>{
                                                            await (0, __TURBOPACK__imported__module__$5b$project$5d2f$EndPoints$2f$Logout$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["Logout"])();
                                                            setIsOpen(false);
                                                            setIsLoggedIn(false);
                                                        },
                                                        className: __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$app$2f$navigation$2f$nav$2e$module$2e$css__$5b$app$2d$client$5d$__$28$css__module$29$__["default"].dropdownItem,
                                                        children: "Logout"
                                                    }, void 0, false, {
                                                        fileName: "[project]/src/app/navigation/page.jsx",
                                                        lineNumber: 95,
                                                        columnNumber: 23
                                                    }, this)
                                                ]
                                            }, void 0, true, {
                                                fileName: "[project]/src/app/navigation/page.jsx",
                                                lineNumber: 93,
                                                columnNumber: 21
                                            }, this)
                                        ]
                                    }, void 0, true, {
                                        fileName: "[project]/src/app/navigation/page.jsx",
                                        lineNumber: 83,
                                        columnNumber: 17
                                    }, this)
                                }, void 0, false, {
                                    fileName: "[project]/src/app/navigation/page.jsx",
                                    lineNumber: 82,
                                    columnNumber: 15
                                }, this)
                            ]
                        }, void 0, true) : /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])(__TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["Fragment"], {
                            children: [
                                /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])(__TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$client$2f$app$2d$dir$2f$link$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["default"], {
                                    className: `${__TURBOPACK__imported__module__$5b$project$5d2f$src$2f$app$2f$navigation$2f$nav$2e$module$2e$css__$5b$app$2d$client$5d$__$28$css__module$29$__["default"].linkWithIcon} ${pathname === "/auth/login" ? __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$app$2f$navigation$2f$nav$2e$module$2e$css__$5b$app$2d$client$5d$__$28$css__module$29$__["default"].active : ""}`,
                                    href: "/auth/login",
                                    onClick: ()=>setIsOpen(false),
                                    children: "Login"
                                }, void 0, false, {
                                    fileName: "[project]/src/app/navigation/page.jsx",
                                    lineNumber: 112,
                                    columnNumber: 15
                                }, this),
                                /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])(__TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$client$2f$app$2d$dir$2f$link$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["default"], {
                                    className: `${__TURBOPACK__imported__module__$5b$project$5d2f$src$2f$app$2f$navigation$2f$nav$2e$module$2e$css__$5b$app$2d$client$5d$__$28$css__module$29$__["default"].linkWithIcon} ${pathname === "/auth/register" ? __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$app$2f$navigation$2f$nav$2e$module$2e$css__$5b$app$2d$client$5d$__$28$css__module$29$__["default"].active : ""}`,
                                    href: "/auth/register",
                                    onClick: ()=>setIsOpen(false),
                                    children: "Register"
                                }, void 0, false, {
                                    fileName: "[project]/src/app/navigation/page.jsx",
                                    lineNumber: 113,
                                    columnNumber: 15
                                }, this)
                            ]
                        }, void 0, true)
                    }, void 0, false, {
                        fileName: "[project]/src/app/navigation/page.jsx",
                        lineNumber: 62,
                        columnNumber: 9
                    }, this)
                ]
            }, void 0, true, {
                fileName: "[project]/src/app/navigation/page.jsx",
                lineNumber: 43,
                columnNumber: 7
            }, this),
            /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])("div", {
                className: `${isLoggedIn ? __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$app$2f$navigation$2f$nav$2e$module$2e$css__$5b$app$2d$client$5d$__$28$css__module$29$__["default"].bottomNav : __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$app$2f$navigation$2f$nav$2e$module$2e$css__$5b$app$2d$client$5d$__$28$css__module$29$__["default"].logged}`,
                children: isLoggedIn && /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])(__TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["Fragment"], {
                    children: [
                        /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])(NavLink, {
                            href: "/",
                            icon: "home",
                            pathname: pathname
                        }, void 0, false, {
                            fileName: "[project]/src/app/navigation/page.jsx",
                            lineNumber: 123,
                            columnNumber: 13
                        }, this),
                        /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])(NavLink, {
                            href: "/posts",
                            icon: "posts",
                            pathname: pathname
                        }, void 0, false, {
                            fileName: "[project]/src/app/navigation/page.jsx",
                            lineNumber: 124,
                            columnNumber: 13
                        }, this),
                        /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])(NavLink, {
                            href: "/groups",
                            icon: "groupe",
                            pathname: pathname
                        }, void 0, false, {
                            fileName: "[project]/src/app/navigation/page.jsx",
                            lineNumber: 125,
                            columnNumber: 13
                        }, this),
                        /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])(NavLink, {
                            href: "/chat",
                            icon: "messages",
                            pathname: pathname
                        }, void 0, false, {
                            fileName: "[project]/src/app/navigation/page.jsx",
                            lineNumber: 126,
                            columnNumber: 13
                        }, this)
                    ]
                }, void 0, true)
            }, void 0, false, {
                fileName: "[project]/src/app/navigation/page.jsx",
                lineNumber: 120,
                columnNumber: 7
            }, this)
        ]
    }, void 0, true, {
        fileName: "[project]/src/app/navigation/page.jsx",
        lineNumber: 42,
        columnNumber: 5
    }, this);
}
_s(Routing, "pCAI/TATZ9cBQD5Q2cxzEs8laLk=", false, function() {
    return [
        __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$navigation$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["usePathname"]
    ];
});
_c = Routing;
// 🔧 Helper for nav links
function NavLink({ href, icon, pathname }) {
    return /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])(__TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$client$2f$app$2d$dir$2f$link$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["default"], {
        className: `${__TURBOPACK__imported__module__$5b$project$5d2f$src$2f$app$2f$navigation$2f$nav$2e$module$2e$css__$5b$app$2d$client$5d$__$28$css__module$29$__["default"].linkWithIcon} ${pathname === href ? __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$app$2f$navigation$2f$nav$2e$module$2e$css__$5b$app$2d$client$5d$__$28$css__module$29$__["default"].active : ""}`,
        href: href,
        children: /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])("span", {
            className: __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$app$2f$navigation$2f$nav$2e$module$2e$css__$5b$app$2d$client$5d$__$28$css__module$29$__["default"].iconWrapper,
            children: [
                /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])(__TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$image$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["default"], {
                    src: `/${icon}2.svg`,
                    alt: icon,
                    width: 25,
                    height: 25,
                    className: __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$app$2f$navigation$2f$nav$2e$module$2e$css__$5b$app$2d$client$5d$__$28$css__module$29$__["default"].iconDefault
                }, void 0, false, {
                    fileName: "[project]/src/app/navigation/page.jsx",
                    lineNumber: 139,
                    columnNumber: 9
                }, this),
                /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$dist$2f$compiled$2f$react$2f$jsx$2d$dev$2d$runtime$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["jsxDEV"])(__TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f$next$2f$image$2e$js__$5b$app$2d$client$5d$__$28$ecmascript$29$__["default"], {
                    src: `/${icon}.svg`,
                    alt: `${icon}-hover`,
                    width: 25,
                    height: 25,
                    className: __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$app$2f$navigation$2f$nav$2e$module$2e$css__$5b$app$2d$client$5d$__$28$css__module$29$__["default"].iconHover
                }, void 0, false, {
                    fileName: "[project]/src/app/navigation/page.jsx",
                    lineNumber: 140,
                    columnNumber: 9
                }, this)
            ]
        }, void 0, true, {
            fileName: "[project]/src/app/navigation/page.jsx",
            lineNumber: 138,
            columnNumber: 7
        }, this)
    }, void 0, false, {
        fileName: "[project]/src/app/navigation/page.jsx",
        lineNumber: 137,
        columnNumber: 5
    }, this);
}
_c1 = NavLink;
var _c, _c1;
__turbopack_context__.k.register(_c, "Routing");
__turbopack_context__.k.register(_c1, "NavLink");
if (typeof globalThis.$RefreshHelpers$ === 'object' && globalThis.$RefreshHelpers !== null) {
    __turbopack_context__.k.registerExports(module, globalThis.$RefreshHelpers$);
}
}}),
}]);

//# sourceMappingURL=_45d153a0._.js.map