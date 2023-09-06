// @ts-ignore
import * as htmlToImage from "html-to-image";
import { toDateForFilename } from "src/lib/util";
import {
  AutoScreenshot,
  LogError,
  ManualScreenshot,
} from "wailsjs/go/main/App";
import type { domain } from "wailsjs/go/models";

export class Screenshot {
  private isFirst: boolean = true;

  private async getScreenshotBase64(
    meta: domain.Meta,
  ): Promise<[string, string]> {
    // Workaround: first screenshot cann't draw values in table.
    const mainPageElem = document.getElementById("stats");
    if (!mainPageElem) {
      throw Error("cann't get element for screenshot");
    }

    if (this.isFirst) {
      await htmlToImage.toPng(mainPageElem);
      this.isFirst = false;
    }
    const dataUrl = await htmlToImage.toPng(mainPageElem);
    const date = toDateForFilename(meta.unixtime);
    const ownShip = meta.own_ship.replaceAll(" ", "-");
    const filename = `${date}_${ownShip}_${meta.arena}_${meta.type}.png`;
    const base64Data = dataUrl.split(",")[1];
    return [filename, base64Data];
  }

  async manual(meta: domain.Meta): Promise<boolean> {
    try {
      const [filename, data] = await this.getScreenshotBase64(meta);
      await ManualScreenshot(filename, data);
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

  async auto(meta: domain.Meta): Promise<boolean> {
    try {
      const [filename, data] = await this.getScreenshotBase64(meta);
      await AutoScreenshot(filename, data);
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
