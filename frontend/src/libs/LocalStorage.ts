import {
  DEFAULT_COLUMN_SETTINGS,
  DEFAULT_PLAYER_NAME_COLUMN_SETTING,
  DEFAULT_SHIP_INFO_COLUMN_SETTING,
} from "./constants";
import type {
  ColumnSettings,
  PlayerNameColumnSetting,
  ShipInfoColumnSetting,
  StatsExtra,
} from "./types";

const KEY_ZOOM_RATE = "zoom_rate";
const KEY_STATS_EXTRA = "stats_extra";
const KEY_PLAYER_NAME_COLUMN_SETTING = "player_name_column_setting";
const KEY_SHIP_INFO_COLUMN_SETTING = "ship_info_column_setting";
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

  getPlayerNameColumnSettings(): PlayerNameColumnSetting {
    const setting = localStorage.getItem(KEY_PLAYER_NAME_COLUMN_SETTING);
    if (!setting) {
      return DEFAULT_PLAYER_NAME_COLUMN_SETTING;
    }

    return JSON.parse(setting) as PlayerNameColumnSetting;
  }

  setPlayerNameColumnSettings(setting: PlayerNameColumnSetting) {
    localStorage.setItem(
      KEY_PLAYER_NAME_COLUMN_SETTING,
      JSON.stringify(setting),
    );
  }

  getShipInfoColumnSettings(): ShipInfoColumnSetting {
    const setting = localStorage.getItem(KEY_SHIP_INFO_COLUMN_SETTING);
    if (!setting) {
      return DEFAULT_SHIP_INFO_COLUMN_SETTING;
    }

    return JSON.parse(setting) as ShipInfoColumnSetting;
  }

  setShipInfoColumnSettings(setting: ShipInfoColumnSetting) {
    localStorage.setItem(KEY_SHIP_INFO_COLUMN_SETTING, JSON.stringify(setting));
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
