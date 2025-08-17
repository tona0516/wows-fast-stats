import { PlayerNameColor } from "./enums";
import type { ColumnSetting, StatsExtra, StatsKey } from "./types";

const KEY_ZOOM_RATE = "zoom_rate";
const KEY_STATS_EXTRA = "stats_extra";
const KEY_PLAYER_NAME_COLOR = "player_name_color";
const KEY_SHOW_CLAN_NATION = "show_clan_nation";
const KEYPREFIX_COLUMN_SETTING = "column_setting_";

export class LocalStorage {
  private static _instance: LocalStorage;
  private constructor() {}
  static get instance(): LocalStorage {
    if (!LocalStorage._instance) {
      LocalStorage._instance = new LocalStorage();
    }
    return LocalStorage._instance;
  }

  getZoomRate(): number {
    const zoomRate = localStorage.getItem(KEY_ZOOM_RATE);
    return zoomRate ? parseInt(zoomRate) : 100;
  }

  setZoomRate(zoomRate: number): void {
    localStorage.setItem(KEY_ZOOM_RATE, zoomRate.toFixed(0));
  }

  getStatsExtra(): StatsExtra {
    const statsExtra = localStorage.getItem(KEY_STATS_EXTRA);
    if (!statsExtra) {
      return "pvp_all";
    }
    return statsExtra as StatsExtra;
  }

  setStatsExtra(statsExtra: StatsExtra): void {
    localStorage.setItem(KEY_STATS_EXTRA, statsExtra);
  }

  getTheme(): string {
    return localStorage.getItem("theme") || "light";
  }

  getPlayerNameColor(): PlayerNameColor {
    const playerNameColor = localStorage.getItem(KEY_PLAYER_NAME_COLOR);
    if (!playerNameColor) {
      return PlayerNameColor.NONE;
    }
    return playerNameColor as PlayerNameColor;
  }

  setPlayerNameColor(playerNameColor: PlayerNameColor) {
    localStorage.setItem(KEY_PLAYER_NAME_COLOR, playerNameColor);
  }

  getShowClanNation(): boolean {
    const showClanNation = localStorage.getItem(KEY_SHOW_CLAN_NATION);
    return showClanNation === "1";
  }

  setShowClanNation(showClanNation: boolean): void {
    localStorage.setItem(KEY_SHOW_CLAN_NATION, showClanNation ? "1" : "0");
  }

  getColumnSetting(key: StatsKey): ColumnSetting {
    const columnSetting = localStorage.getItem(KEYPREFIX_COLUMN_SETTING + key);
    if (!columnSetting) {
      return { ship: false, overall: false, digit: 0 };
    }
    return JSON.parse(columnSetting) as ColumnSetting;
  }

  setColumnSetting(key: StatsKey, columnSetting: Partial<ColumnSetting>) {
    const currentSetting = this.getColumnSetting(key);
    const newSetting: ColumnSetting = {
      ...currentSetting,
      ...columnSetting,
    };
    localStorage.setItem(
      KEYPREFIX_COLUMN_SETTING + key,
      JSON.stringify(newSetting),
    );
  }
}
