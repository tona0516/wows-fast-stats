import { format, fromUnixTime } from "date-fns";
import * as htmlToImage from "html-to-image";
import {
  AutoScreenshot,
  LogError,
  ManualScreenshot,
} from "wailsjs/go/main/App";
import type { domain } from "wailsjs/go/models";

export class Screenshot {
  private isFirst: boolean = true;

  constructor(private targetElementID: string) {}

  async manual(meta: domain.Meta): Promise<boolean> {
    const filename = deriveFileName(meta);
    const image = await this.getBase64Image(meta);
    return await ManualScreenshot(filename, image);
  }

  async auto(meta: domain.Meta) {
    const filename = deriveFileName(meta);
    const image = await this.getBase64Image(meta);
    await AutoScreenshot(filename, image);
  }

  private async getBase64Image(meta: domain.Meta): Promise<string> {
    try {
      const element = document.getElementById(this.targetElementID);

      // Workaround: first screenshot can't draw values in table.
      if (this.isFirst) {
        await htmlToImage.toPng(element!);
        this.isFirst = false;
      }
      return (await htmlToImage.toPng(element!)).split(",")[1];
    } catch (error) {
      LogError(JSON.stringify(error));
      throw error;
    }
  }
}

const deriveFileName = (meta: domain.Meta): string =>
  `${format(
    fromUnixTime(meta.unixtime),
    "yyyy-MM-dd-HH-mm-ss",
  )}_${meta.own_ship.replaceAll(" ", "-")}_${meta.arena}_${meta.type}.png`;
