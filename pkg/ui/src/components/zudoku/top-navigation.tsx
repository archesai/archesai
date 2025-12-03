// import { cx } from "class-variance-authority";
// import { deepEqual } from "fast-equals";
// import { Suspense } from "react";
// import type { NavLinkProps } from "react-router";
// import { NavLink } from "react-router";
// import type { NavigationItem } from "../../config/validators/NavigationSchema.js";
// import { useAuth } from "../authentication/hook.js";
// import { joinUrl } from "../util/joinUrl.js";
// import { useCurrentNavigation, useZudoku } from "./context/ZudokuContext.js";
// import { shouldShowItem, traverseNavigationItem } from "./navigation/utils.js";
// import { Slot } from "./Slot.js";

// export const TopNavigation = () => {
//   const context = useZudoku();
//   const { navigation } = context;
//   const auth = useAuth();
//   const filteredItems = navigation.filter(shouldShowItem(auth, context));

//   if (filteredItems.length === 0 || import.meta.env.MODE === "standalone") {
//     return <style>{`:root { --top-nav-height: 0px; }`}</style>;
//   }

//   return (
//     <Suspense>
//       <div className="relative hidden h-(--top-nav-height) items-center justify-between px-8 text-sm lg:flex">
//         <nav className="text-sm">
//           <ul className="flex flex-row items-center gap-8">
//             {filteredItems.map((item) => (
//               <li key={item.label + item.type}>
//                 <TopNavItem {...item} />
//               </li>
//             ))}
//           </ul>
//         </nav>
//         <Slot.Target name="top-navigation-side" />
//       </div>
//       {/* <PageProgress /> */}
//     </Suspense>
//   );
// };

// const getPathForItem = (item: NavigationItem): string => {
//   switch (item.type) {
//     case "doc":
//       return joinUrl(item.path);
//     case "link":
//       return item.to;
//     case "category": {
//       if (item.link?.path) {
//         return joinUrl(item.link.path);
//       }

//       return (
//         traverseNavigationItem(item, (child) => {
//           if (child.type !== "category") {
//             return getPathForItem(child);
//           }
//         }) ?? ""
//       );
//     }
//     case "custom-page":
//       return item.path;
//   }
// };

// export const TopNavLink = ({
//   isActive,
//   children,
//   ...props
// }: {
//   isActive?: boolean;
//   children: React.ReactNode;
// } & NavLinkProps) => {
//   return (
//     <NavLink
//       className={({ isActive: isActiveNavLink, isPending }) => {
//         const isActiveReal = isActiveNavLink || isActive;
//         return cx(
//           "-mb-px relative flex items-center gap-2 font-medium transition delay-75 duration-150 lg:py-3.5",
//           isActiveReal || isPending
//             ? [
//                 "text-foreground",
//                 // underline with view transition animation
//                 "after:absolute after:right-0 after:bottom-0 after:left-0 after:content-['']",
//                 "after:h-0.5 after:bg-primary",
//                 isActiveReal &&
//                   "after:[view-transition-name:top-nav-underline]",
//                 isPending && "after:bg-primary/25",
//               ]
//             : "text-foreground/75 hover:text-foreground",
//         );
//       }}
//       viewTransition
//       {...props}
//     >
//       {children}
//     </NavLink>
//   );
// };

// export const TopNavItem = (item: NavigationItem) => {
//   const currentNav = useCurrentNavigation();
//   const isActiveTopNavItem = deepEqual(currentNav.topNavItem, item);

//   const path = getPathForItem(item);

//   return (
//     // We don't use isActive here because it has to be inside the navigation,
//     // the top nav id doesn't necessarily start with the navigation id
//     <TopNavLink
//       isActive={isActiveTopNavItem}
//       to={path}
//     >
//       {item.icon && (
//         <item.icon
//           className="align-[-0.125em]"
//           size={16}
//         />
//       )}
//       {item.label}
//     </TopNavLink>
//   );
// };
