import type { IColumn } from "src/lib/column/intetface/IColumn";
import type { BasicKey } from "src/lib/types";
import ShipInfoTableData from "src/tabledata_component/ShipInfoTableData.svelte";
import type { domain } from "wailsjs/go/models";

import CommonWealthNationIcon from "src/assets/images/nation-commonwealth.png";
import EuropeNationIcon from "src/assets/images/nation-europe.png";
import FranceNationIcon from "src/assets/images/nation-france.png";
import GermanyNationIcon from "src/assets/images/nation-germany.png";
import ItalyNationIcon from "src/assets/images/nation-italy.png";
import JapanNationIcon from "src/assets/images/nation-japan.png";
import NetherlandsNationIcon from "src/assets/images/nation-netherlands.png";
import NoneNationIcon from "src/assets/images/nation-none.png";
import PanAmericaNationIcon from "src/assets/images/nation-pan-america.png";
import PanAsiaNationIcon from "src/assets/images/nation-pan-asia.png";
import SpainNationIcon from "src/assets/images/nation-spain.png";
import UkNationIcon from "src/assets/images/nation-uk.png";
import UsaNationIcon from "src/assets/images/nation-usa.png";
import UssrNationIcon from "src/assets/images/nation-ussr.png";

import PremiumBBShipIcon from "src/assets/images/icon-bb-premium.png";
import BBShipIcon from "src/assets/images/icon-bb.png";
import PremiumCLShipIcon from "src/assets/images/icon-cl-premium.png";
import CLShipIcon from "src/assets/images/icon-cl.png";
import PremiumCVShipIcon from "src/assets/images/icon-cv-premium.png";
import CVShipIcon from "src/assets/images/icon-cv.png";
import PremiumDDShipIcon from "src/assets/images/icon-dd-premium.png";
import DDShipIcon from "src/assets/images/icon-dd.png";
import NoneShipIcon from "src/assets/images/icon-none.png";
import PremiumSSShipIcon from "src/assets/images/icon-ss-premium.png";
import SSShipIcon from "src/assets/images/icon-ss.png";
import { BASE_NUMBERS_URL } from "src/const";

const NATION_ICONS: { [key: string]: string } = {
  japan: JapanNationIcon,
  usa: UsaNationIcon,
  ussr: UssrNationIcon,
  germany: GermanyNationIcon,
  uk: UkNationIcon,
  france: FranceNationIcon,
  italy: ItalyNationIcon,
  pan_asia: PanAsiaNationIcon,
  europe: EuropeNationIcon,
  netherlands: NetherlandsNationIcon,
  commonwealth: CommonWealthNationIcon,
  pan_america: PanAmericaNationIcon,
  spain: SpainNationIcon,
  none: NoneNationIcon,
};

export class ShipInfo implements IColumn<BasicKey> {
  constructor(private userConfig: domain.UserConfig) {}

  displayKey(): BasicKey {
    return "ship_info";
  }

  minDisplayName(): string {
    return "艦";
  }

  fullDisplayName(): string {
    return "艦情報";
  }

  shouldShowColumn(): boolean {
    return true;
  }

  countInnerColumn(): number {
    return 3;
  }

  svelteComponent() {
    return ShipInfoTableData;
  }

  displayValue(player: domain.Player): string {
    return `${this.tierString(player.ship_info.tier)} ${player.ship_info.name}`;
  }

  bgColorCode(player: domain.Player): string {
    const ownColor = this.userConfig.custom_color.ship_type.own;
    switch (player.ship_info.type) {
      case "cv":
        return ownColor.cv;
      case "bb":
        return ownColor.bb;
      case "cl":
        return ownColor.cl;
      case "dd":
        return ownColor.dd;
      case "ss":
        return ownColor.ss;
      default:
        return "#00000000";
    }
  }

  shipURL(player: domain.Player): string {
    return (
      BASE_NUMBERS_URL +
      "ship/" +
      player.ship_info.id +
      "," +
      player.ship_info.name.replaceAll(" ", "-")
    );
  }

  shipTypeIconPath(player: domain.Player): string {
    const shipInfo = player.ship_info;

    switch (shipInfo.type) {
      case "cv":
        return shipInfo.is_premium ? PremiumCVShipIcon : CVShipIcon;
      case "bb":
        return shipInfo.is_premium ? PremiumBBShipIcon : BBShipIcon;
      case "cl":
        return shipInfo.is_premium ? PremiumCLShipIcon : CLShipIcon;
      case "dd":
        return shipInfo.is_premium ? PremiumDDShipIcon : DDShipIcon;
      case "ss":
        return shipInfo.is_premium ? PremiumSSShipIcon : SSShipIcon;
      default:
        return NoneShipIcon;
    }
  }

  nationIconPath(player: domain.Player): string {
    return NATION_ICONS[player.ship_info.nation] ?? NATION_ICONS.none;
  }

  private tierString(value: number): string {
    if (value === 11) return "★";

    const decimal = [10, 9, 5, 4, 1];
    const romanNumeral = ["X", "IX", "V", "IV", "I"];

    let romanized = "";

    for (var i = 0; i < decimal.length; i++) {
      while (decimal[i] <= value) {
        romanized += romanNumeral[i];
        value -= decimal[i];
      }
    }

    return romanized;
  }
}
