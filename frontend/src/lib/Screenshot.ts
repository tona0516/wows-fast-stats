import { format, fromUnixTime } from "date-fns";
import * as htmlToImage from "html-to-image";
import {
  AutoScreenshot,
  LogError,
  ManualScreenshot,
} from "wailsjs/go/main/App";
import type { model } from "wailsjs/go/models";

export class Screenshot {
  private isFirst: boolean = true;

  constructor(private targetElementID: string) {}

  async manual(meta: model.Meta): Promise<boolean> {
    const filename = deriveFileName(meta);
    const image = await this.getBase64Image();
    return await ManualScreenshot(filename, image);
  }

  async auto(meta: model.Meta) {
    const filename = deriveFileName(meta);
    const image = await this.getBase64Image();
    await AutoScreenshot(filename, image);
  }

  private async getBase64Image(): Promise<string> {
    try {
      const element = document.getElementById(this.targetElementID);

      // Workaround: first screenshot can't draw values in table.
      if (this.isFirst) {
        await htmlToImage.toPng(element!);
        this.isFirst = false;
      }
      return (await htmlToImage.toPng(element!)).split(",")[1];
    } catch (error) {
      LogError(`${this.getBase64Image.name}: ${JSON.stringify(error)}`, {});
      throw error;
    }
  }
}

const deriveFileName = (meta: model.Meta): string => {
  const items = [
    format(fromUnixTime(meta.unixtime), "yyyy-MM-dd-HH-mm-ss"),
    meta.own_ship.replaceAll(" ", "-"),
    meta.arena,
    meta.type,
  ];

  return `${items.join("_")}.png`;
};
