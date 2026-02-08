import type { PlayerNameColorType, StatsExtra, StatsKey } from "./types";

export interface DisplayPref {
  version: number;
  zoomRate: number;
  statsExtra: StatsExtra;
  showBoarder: boolean;
  showSiPrefix: boolean;
  player: {
    enableNationFlag: boolean;
    colorType: PlayerNameColorType;
  };
  warship: {
    enableNationFlag: boolean;
  };
  // ship/overall columns
  pr: {
    showShip: boolean;
    showOverall: boolean;
    digit: number;
  };
  winRate: {
    showShip: boolean;
    showOverall: boolean;
    digit: number;
  };
  damage: {
    showShip: boolean;
    showOverall: boolean;
    digit: number;
  };
  maxDamage: {
    showShip: boolean;
    showOverall: boolean;
  };
  kdRate: {
    showShip: boolean;
    showOverall: boolean;
    digit: number;
  };
  kill: {
    showShip: boolean;
    showOverall: boolean;
    digit: number;
  };
  exp: {
    showShip: boolean;
    showOverall: boolean;
    digit: number;
  };
  battles: {
    showShip: boolean;
    showOverall: boolean;
  };
  platoonRate: {
    showShip: boolean;
    showOverall: boolean;
    digit: number;
  };
  efficiencyBadge: {
    showShip: boolean;
    showOverall: boolean;
  };
  // ship only columns
  planesKilled: {
    showShip: boolean;
    digit: number;
  };
  survivedRate: {
    showShip: boolean;
    digit: number;
  };
  hitRate: {
    showShip: boolean;
    digit: number;
  };
  // overall only columns
  avgTier: {
    showOverall: boolean;
    digit: number;
  };
  threatLevel: {
    showOverall: boolean;
    digit: number;
  };
  usingShipTypeRate: {
    showOverall: boolean;
    digit: number;
  };
  usingTierRate: {
    showOverall: boolean;
    digit: number;
  };
  columnOrder: {
    ship: StatsKey[];
    overall: StatsKey[];
  };
}

export const DEFAULT_DISPLAY_PREF: DisplayPref = {
  version: 1,
  zoomRate: 100,
  statsExtra: "pvpAll",
  showBoarder: false,
  showSiPrefix: false,
  player: {
    enableNationFlag: true,
    colorType: "none",
  },
  warship: {
    enableNationFlag: true,
  },
  pr: {
    showShip: true,
    showOverall: true,
    digit: 0,
  },
  winRate: {
    showShip: true,
    showOverall: true,
    digit: 1,
  },
  damage: {
    showShip: true,
    showOverall: false,
    digit: 0,
  },
  maxDamage: {
    showShip: false,
    showOverall: false,
  },
  kdRate: {
    showShip: true,
    showOverall: false,
    digit: 2,
  },
  kill: {
    showShip: false,
    showOverall: false,
    digit: 2,
  },
  exp: {
    showShip: false,
    showOverall: false,
    digit: 0,
  },
  battles: {
    showShip: true,
    showOverall: true,
  },
  platoonRate: {
    showShip: false,
    showOverall: false,
    digit: 2,
  },
  efficiencyBadge: {
    showShip: true,
    showOverall: false,
  },
  planesKilled: {
    showShip: false,
    digit: 1,
  },
  survivedRate: {
    showShip: false,
    digit: 1,
  },
  hitRate: {
    showShip: false,
    digit: 1,
  },
  avgTier: {
    showOverall: true,
    digit: 2,
  },
  threatLevel: {
    showOverall: false,
    digit: 0,
  },
  usingShipTypeRate: {
    showOverall: false,
    digit: 1,
  },
  usingTierRate: {
    showOverall: false,
    digit: 1,
  },
  columnOrder: {
    ship: [
      "pr",
      "maxDamage",
      "damage",
      "winRate",
      "kdRate",
      "kill",
      "exp",
      "battles",
      "survivedRate",
      "platoonRate",
      "efficiencyBadge",
      "planesKilled",
      "hitRate",
    ],
    overall: [
      "pr",
      "winRate",
      "damage",
      "maxDamage",
      "kdRate",
      "kill",
      "exp",
      "battles",
      "survivedRate",
      "platoonRate",
      "efficiencyBadge",
      "threatLevel",
      "avgTier",
      "usingTierRate",
      "usingShipTypeRate",
    ],
  },
};
