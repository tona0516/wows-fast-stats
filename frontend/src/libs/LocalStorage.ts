import { DEFAULT_COLUMN_SETTINGS } from "./constants";
import type { ColumnSettings, PlayerNameColor, StatsExtra } from "./types";

const KEY_ZOOM_RATE = "zoom_rate";
const KEY_STATS_EXTRA = "stats_extra";
const KEY_PLAYER_NAME_COLOR = "player_name_color";
const KEY_SHOW_CLAN_NATION = "show_clan_nation";
const KEY_COLUMN_SETTINGS = "column_settings";

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

  getPlayerNameColor(): PlayerNameColor {
    const playerNameColor = localStorage.getItem(KEY_PLAYER_NAME_COLOR);
    if (!playerNameColor) {
      return "none";
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

  getColumnSettings(): ColumnSettings {
    const columnSettings = localStorage.getItem(KEY_COLUMN_SETTINGS);
    if (!columnSettings) {
      return DEFAULT_COLUMN_SETTINGS;
    }
    return JSON.parse(columnSettings) as ColumnSettings;
  }

  setColumnSettings(settings: ColumnSettings) {
    localStorage.setItem(KEY_COLUMN_SETTINGS, JSON.stringify(settings));
  }
}
