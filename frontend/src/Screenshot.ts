import * as htmlToImage from "html-to-image";
import { toDateForFilename } from "src/util";
import { AutoScreenshot, ManualScreenshot } from "wailsjs/go/main/App";
import type { domain } from "wailsjs/go/models";

export class Screenshot {
  isFirst: boolean;
  battle: domain.Battle;

  constructor(battle: domain.Battle, isFirst: boolean) {
    this.battle = battle;
    this.isFirst = isFirst;
  }

  private async getScreenshotBase64(): Promise<[string, string]> {
    // Workaround: first screenshot cann't draw values in table.
    const mainPageElem = document.getElementById("stats");
    if (this.isFirst) {
      await htmlToImage.toPng(mainPageElem);
    }
    const dataUrl = await htmlToImage.toPng(mainPageElem);
    const date = toDateForFilename(this.battle.meta.unixtime);
    const ownShip = this.battle.meta.own_ship.replaceAll(" ", "-");
    const filename = `${date}_${ownShip}_${this.battle.meta.arena}_${this.battle.meta.type}.png`;
    const base64Data = dataUrl.split(",")[1];
    return [filename, base64Data];
  }

  async manual() {
    const [filename, data] = await this.getScreenshotBase64();
    await ManualScreenshot(filename, data);
  }

  async auto() {
    const [filename, data] = await this.getScreenshotBase64();
    await AutoScreenshot(filename, data);
  }
}
