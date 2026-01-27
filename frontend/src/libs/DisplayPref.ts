import type {
  PlayerNameColorType,
  StatsExtra,
  WarshipNamesColorType,
} from "./types";

export interface DisplayPref {
  version: number;
  zoomRate: number;
  statsExtra: StatsExtra;
  player: {
    enableNationFlag: boolean;
    colorType: PlayerNameColorType;
  };
  warship: {
    enableNationFlag: boolean;
    colorType: WarshipNamesColorType;
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
    kmgUnit: boolean;
  };
  maxDamage: {
    showShip: boolean;
    showOverall: boolean;
    kmgUnit: boolean;
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
    kmgUnit: boolean;
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
    kmgUnit: boolean;
  };
  usingShipTypeRate: {
    showOverall: boolean;
    digit: number;
  };
  usingTierRate: {
    showOverall: boolean;
    digit: number;
  };
}

export const DEFAULT_DISPLAY_PREF: DisplayPref = {
  version: 1,
  zoomRate: 100,
  statsExtra: "pvpAll",
  player: {
    enableNationFlag: true,
    colorType: "none",
  },
  warship: {
    enableNationFlag: true,
    colorType: "none",
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
    kmgUnit: false,
  },
  maxDamage: {
    showShip: false,
    showOverall: false,
    kmgUnit: false,
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
    kmgUnit: false,
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
    kmgUnit: false,
  },
  usingShipTypeRate: {
    showOverall: false,
    digit: 1,
  },
  usingTierRate: {
    showOverall: false,
    digit: 1,
  },
};
