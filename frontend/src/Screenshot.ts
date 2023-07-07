import * as htmlToImage from "html-to-image";
import { AutoScreenshot, ManualScreenshot } from "../wailsjs/go/main/App";
import type { vo } from "../wailsjs/go/models";
import { toDateForFilename } from "./util";

export class Screenshot {
  isFirst: boolean;
  battle: vo.Battle;

  constructor(battle: vo.Battle, isFirst: boolean) {
    this.battle = battle;
    this.isFirst = isFirst;
  }

  private async getScreenshotBase64(): Promise<[string, string]> {
    // Workaround: first screenshot cann't draw values in table.
    const mainPageElem = document.getElementById("mainpage");
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
