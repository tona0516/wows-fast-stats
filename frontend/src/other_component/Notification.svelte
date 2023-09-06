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

  export function showSuccessToast(description: string) {
    toasts.add({
      description: description,
      duration: DURATION_MS,
      placement: PLACEMENT,
      type: "success",
      theme: THEME,
    });
  }

  export function showErrorToast(message: any) {
    toasts.add({
      description: toDescription(message),
      duration: DURATION_MS,
      placement: PLACEMENT,
      type: "error",
      theme: THEME,
    });
  }

  export function showToastWithKey(
    message: any,
    type: ToastType,
    key: string,
    onClick?: () => void,
  ) {
    if (toastDict[key]) {
      return;
    }

    if (!onClick) {
      onClick = () => {};
    }

    const toast = toasts.add({
      description: toDescription(message),
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
    const toast = toastDict[key];

    if (!toast || !toast.remove) {
      return;
    }

    toast.remove();
    delete toastDict[key];
  }

  const toDescription = (message: any): string => {
    return message instanceof Error
      ? `${message.name}: ${message.message}`
      : message;
  };
</script>

<ToastContainer let:data>
  <FlatToast {data} />
</ToastContainer>
