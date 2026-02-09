import type { PlayerNameColorType, StatsExtra, StatsKey } from "./types";

export interface DisplayPref {
  version: number;
  zoomRate: number;
  statsExtra: StatsExtra;
  isBorderVisible: boolean;
  isSiPrefixEnabled: boolean;
  player: {
    isNationFlagEnabled: boolean;
    colorType: PlayerNameColorType;
  };
  warship: {
    isNationFlagEnabled: boolean;
  };
  // ship/overall columns
  pr: {
    isShipVisible: boolean;
    isOverallVisible: boolean;
    digit: number;
  };
  winRate: {
    isShipVisible: boolean;
    isOverallVisible: boolean;
    digit: number;
  };
  damage: {
    isShipVisible: boolean;
    isOverallVisible: boolean;
    digit: number;
  };
  maxDamage: {
    isShipVisible: boolean;
    isOverallVisible: boolean;
  };
  kdRate: {
    isShipVisible: boolean;
    isOverallVisible: boolean;
    digit: number;
  };
  kill: {
    isShipVisible: boolean;
    isOverallVisible: boolean;
    digit: number;
  };
  exp: {
    isShipVisible: boolean;
    isOverallVisible: boolean;
    digit: number;
  };
  battles: {
    isShipVisible: boolean;
    isOverallVisible: boolean;
  };
  platoonRate: {
    isShipVisible: boolean;
    isOverallVisible: boolean;
    digit: number;
  };
  efficiencyBadge: {
    isShipVisible: boolean;
    isOverallVisible: boolean;
  };
  // ship only columns
  planesKilled: {
    isShipVisible: boolean;
    digit: number;
  };
  survivedRate: {
    isShipVisible: boolean;
    digit: number;
  };
  hitRate: {
    isShipVisible: boolean;
    digit: number;
  };
  // overall only columns
  avgTier: {
    isOverallVisible: boolean;
    digit: number;
  };
  usingShipTypeRate: {
    isOverallVisible: boolean;
    digit: number;
  };
  usingTierRate: {
    isOverallVisible: boolean;
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
  isBorderVisible: false,
  isSiPrefixEnabled: false,
  player: {
    isNationFlagEnabled: true,
    colorType: "none",
  },
  warship: {
    isNationFlagEnabled: true,
  },
  pr: {
    isShipVisible: true,
    isOverallVisible: true,
    digit: 0,
  },
  winRate: {
    isShipVisible: true,
    isOverallVisible: true,
    digit: 1,
  },
  damage: {
    isShipVisible: true,
    isOverallVisible: false,
    digit: 0,
  },
  maxDamage: {
    isShipVisible: false,
    isOverallVisible: false,
  },
  kdRate: {
    isShipVisible: true,
    isOverallVisible: false,
    digit: 2,
  },
  kill: {
    isShipVisible: false,
    isOverallVisible: false,
    digit: 2,
  },
  exp: {
    isShipVisible: false,
    isOverallVisible: false,
    digit: 0,
  },
  battles: {
    isShipVisible: true,
    isOverallVisible: true,
  },
  platoonRate: {
    isShipVisible: false,
    isOverallVisible: false,
    digit: 2,
  },
  efficiencyBadge: {
    isShipVisible: true,
    isOverallVisible: false,
  },
  planesKilled: {
    isShipVisible: false,
    digit: 1,
  },
  survivedRate: {
    isShipVisible: false,
    digit: 1,
  },
  hitRate: {
    isShipVisible: false,
    digit: 1,
  },
  avgTier: {
    isOverallVisible: true,
    digit: 2,
  },
  usingShipTypeRate: {
    isOverallVisible: false,
    digit: 1,
  },
  usingTierRate: {
    isOverallVisible: false,
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
      "avgTier",
      "usingTierRate",
      "usingShipTypeRate",
    ],
  },
};
