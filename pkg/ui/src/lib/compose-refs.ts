import { useCallback } from "react";

type PossibleRef<T> = React.Ref<T> | undefined;

/**
 * Composes multiple refs into a single callback ref.
 * Fully React 19 compatible: supports cleanup functions.
 */
function composeRefs<T>(...refs: PossibleRef<T>[]): React.RefCallback<T> {
  return (node) => {
    const cleanups: ((() => void) | undefined)[] = [];

    for (const ref of refs) {
      let cleanup: (() => void) | undefined;
      if (typeof ref === "function") {
        const result = ref(node);
        cleanup = typeof result === "function" ? result : undefined;
      } else {
        setRef(ref, node);
        cleanup = undefined;
      }
      cleanups.push(cleanup);
    }

    // Return a cleanup function in React 19 style if any ref returned one
    return () => {
      for (let i = 0; i < refs.length; i++) {
        const cleanup = cleanups[i];
        if (typeof cleanup === "function") {
          cleanup();
        } else {
          setRef(refs[i], null);
        }
      }
    };
  };
}

/**
 * A utility to set a ref, whether it's a function or a ref object.
 */
function setRef<T>(ref: PossibleRef<T>, value: null | T): void {
  if (typeof ref === "function") {
    ref(value);
  } else if (ref != null) {
    // ref object (e.g. useRef)
    ref.current = value;
  }
}

/**
 * A hook that memoizes a composed ref.
 */
function useComposedRefs<T>(...refs: PossibleRef<T>[]): React.RefCallback<T> {
  // Create stable callback without dependencies since refs are passed directly
  return useCallback(
    (node: null | T) => {
      return composeRefs(...refs)(node);
    },
    [refs],
  );
}

export { composeRefs, useComposedRefs };
