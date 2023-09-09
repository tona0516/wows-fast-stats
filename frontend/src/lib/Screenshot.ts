// @ts-ignore
import { format, fromUnixTime } from "date-fns";
import * as htmlToImage from "html-to-image";
import { MAIN_PAGE_ID } from "src/const";
import { ScreenshotType } from "src/lib/types";
import {
  AutoScreenshot,
  LogError,
  ManualScreenshot,
} from "wailsjs/go/main/App";
import type { domain } from "wailsjs/go/models";

export class Screenshot {
  private isFirst: boolean = true;

  async take(type: ScreenshotType, meta: domain.Meta): Promise<boolean> {
    try {
      const element = document.getElementById(MAIN_PAGE_ID);

      // Workaround: first screenshot can't draw values in table.
      if (this.isFirst) {
        await htmlToImage.toPng(element!);
        this.isFirst = false;
      }
      const base64Data = (await htmlToImage.toPng(element!)).split(",")[1];

      const filename = `${format(
        fromUnixTime(meta.unixtime),
        "yyyy-MM-dd-HH-mm-ss",
      )}_${meta.own_ship.replaceAll(" ", "-")}_${meta.arena}_${meta.type}.png`;

      switch (type) {
        case ScreenshotType.manual:
          await ManualScreenshot(filename, base64Data);
          break;
        case ScreenshotType.auto:
          await AutoScreenshot(filename, base64Data);
          break;
      }

      return true;
    } catch (error) {
      const errorJSON = JSON.stringify(error);
      if (errorJSON.includes("ユーザキャンセル")) {
        return false;
      }
      LogError(errorJSON);
      throw error;
    }
  }
}
