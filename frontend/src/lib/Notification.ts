import UIkit from "uikit";

const POSITION = "top-right";

export class Notification {
  success(message: string, durationMs: number = 3000) {
    UIkit.notification({
      message: this.injectMessage("check", message),
      timeout: durationMs,
      pos: POSITION,
    });
  }

  failure(error: unknown, durationMs: number = 0) {
    let message: string = "";
    if (error instanceof Error) {
      message = `${error.name}: ${error.message}`;
    } else if (typeof error === "string") {
      message = error;
    } else {
      message = JSON.stringify(error);
    }

    UIkit.notification({
      message: this.injectMessage("ban", message),
      timeout: durationMs,
      pos: POSITION,
    });
  }

  private injectMessage(icon: string, message: string): string {
    return `<div class="uk-text-small">
        <UkIcon name=${icon} />
        <span class="uk-text-middle">${message}</span>
    </div>`;
  }
}
