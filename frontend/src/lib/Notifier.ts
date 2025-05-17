import UIkit from "uikit";

const POSITION = "top-right";

export namespace Notifier {
  export const success = (message: string, durationMs = 3000) => {
    notify(message, "check", durationMs);
  };

  export const failure = (error: unknown, durationMs = 0) => {
    let message = "";
    if (error instanceof Error) {
      message = `${error.name}: ${error.message}`;
    } else if (typeof error === "string") {
      message = error;
    } else {
      message = JSON.stringify(error);
    }

    notify(message, "ban", durationMs);
  };
}

const notify = (message: string, icon: string, durationMs: number) => {
  UIkit.notification({
    message: `<div class="uk-text-small">
                <UkIcon name=${icon} />
                <span class="uk-text-middle">${message}</span>
              </div>`,
    timeout: durationMs,
    pos: POSITION,
  });
};
