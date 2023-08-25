<script lang="ts">
  import { toasts, ToastContainer, FlatToast } from "svelte-toasts";
  import type {
    Placement,
    Theme,
    ToastProps,
    ToastType,
  } from "svelte-toasts/types/common";

  const DURATION_MS = 5000;
  const PLACEMENT: Placement = "bottom-right";
  const THEME: Theme = "dark";

  let toastDict: { [key: string]: ToastProps } = {};

  export function showSuccessToast(message: string) {
    toasts.add({
      description: message,
      duration: DURATION_MS,
      placement: PLACEMENT,
      type: "success",
      theme: THEME,
    });
  }

  export function showErrorToast(message: string) {
    toasts.add({
      description: message,
      duration: DURATION_MS,
      placement: PLACEMENT,
      type: "error",
      theme: THEME,
    });
  }

  export function showToastWithKey(
    message: string,
    type: ToastType,
    key: string,
    onClick: () => void = null
  ) {
    if (toastDict[key]) {
      return;
    }

    if (!onClick) {
      onClick = () => {};
    }

    const toast = toasts.add({
      description: message,
      duration: 0,
      placement: PLACEMENT,
      type: type,
      theme: THEME,
      onClick: onClick,
    });
    toastDict[key] = toast;
    return;
  }

  export function removeToastWithKey(key: string) {
    if (!toastDict[key]) {
      return;
    }

    toastDict[key].remove();
    delete toastDict[key];
  }
</script>

<ToastContainer let:data>
  <FlatToast {data} />
</ToastContainer>
