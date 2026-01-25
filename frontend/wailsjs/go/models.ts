export namespace core {
	
	export class TeamThreatLevel {
	    average: number;
	    dissociationDegree: number;
	    accuracy: number;
	
	    static createFrom(source: any = {}) {
	        return new TeamThreatLevel(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.average = source["average"];
	        this.dissociationDegree = source["dissociationDegree"];
	        this.accuracy = source["accuracy"];
	    }
	}
	export class TeamAverageStats {
	    ship_pr: number;
	    ship_damage: number;
	    ship_win_rate: number;
	    ship_battles: number;
	    overall_pr: number;
	    overall_damage: number;
	    overall_win_rate: number;
	    overall_battles: number;
	
	    static createFrom(source: any = {}) {
	        return new TeamAverageStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ship_pr = source["ship_pr"];
	        this.ship_damage = source["ship_damage"];
	        this.ship_win_rate = source["ship_win_rate"];
	        this.ship_battles = source["ship_battles"];
	        this.overall_pr = source["overall_pr"];
	        this.overall_damage = source["overall_damage"];
	        this.overall_win_rate = source["overall_win_rate"];
	        this.overall_battles = source["overall_battles"];
	    }
	}
	export class TeamStats {
	    teamAverageStats: TeamAverageStats;
	    teamThreatLevel: TeamThreatLevel;
	
	    static createFrom(source: any = {}) {
	        return new TeamStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.teamAverageStats = this.convertValues(source["teamAverageStats"], TeamAverageStats);
	        this.teamThreatLevel = this.convertValues(source["teamThreatLevel"], TeamThreatLevel);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class EfficiencyBadgeGroup {
	    expert: number;
	    first: number;
	    second: number;
	    third: number;
	
	    static createFrom(source: any = {}) {
	        return new EfficiencyBadgeGroup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.expert = source["expert"];
	        this.first = source["first"];
	        this.second = source["second"];
	        this.third = source["third"];
	    }
	}
	export class TierGroup {
	    low: number;
	    middle: number;
	    high: number;
	
	    static createFrom(source: any = {}) {
	        return new TierGroup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.low = source["low"];
	        this.middle = source["middle"];
	        this.high = source["high"];
	    }
	}
	export class ShipTypeGroup {
	    cv: number;
	    bb: number;
	    cl: number;
	    dd: number;
	    ss: number;
	
	    static createFrom(source: any = {}) {
	        return new ShipTypeGroup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cv = source["cv"];
	        this.bb = source["bb"];
	        this.cl = source["cl"];
	        this.dd = source["dd"];
	        this.ss = source["ss"];
	    }
	}
	export class ThreatLevel {
	    rank: string;
	    raw: number;
	    modified: number;
	
	    static createFrom(source: any = {}) {
	        return new ThreatLevel(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.rank = source["rank"];
	        this.raw = source["raw"];
	        this.modified = source["modified"];
	    }
	}
	export class OverallStats {
	    battles: number;
	    damage: RatingValue;
	    maxDamage: MaxDamage;
	    winRate: RatingValue;
	    survivedRate: SurvivedRate;
	    kdRate: number;
	    kill: number;
	    exp: number;
	    pr: RatingValue;
	    threatLevel: ThreatLevel;
	    avgTier: number;
	    usingShipTypeRate: ShipTypeGroup;
	    usingTierRate: TierGroup;
	    platoonRate: number;
	    efficiencyBadge: EfficiencyBadgeGroup;
	
	    static createFrom(source: any = {}) {
	        return new OverallStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.battles = source["battles"];
	        this.damage = this.convertValues(source["damage"], RatingValue);
	        this.maxDamage = this.convertValues(source["maxDamage"], MaxDamage);
	        this.winRate = this.convertValues(source["winRate"], RatingValue);
	        this.survivedRate = this.convertValues(source["survivedRate"], SurvivedRate);
	        this.kdRate = source["kdRate"];
	        this.kill = source["kill"];
	        this.exp = source["exp"];
	        this.pr = this.convertValues(source["pr"], RatingValue);
	        this.threatLevel = this.convertValues(source["threatLevel"], ThreatLevel);
	        this.avgTier = source["avgTier"];
	        this.usingShipTypeRate = this.convertValues(source["usingShipTypeRate"], ShipTypeGroup);
	        this.usingTierRate = this.convertValues(source["usingTierRate"], TierGroup);
	        this.platoonRate = source["platoonRate"];
	        this.efficiencyBadge = this.convertValues(source["efficiencyBadge"], EfficiencyBadgeGroup);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class HitRate {
	    mainBattery: number;
	    torpedoes: number;
	
	    static createFrom(source: any = {}) {
	        return new HitRate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mainBattery = source["mainBattery"];
	        this.torpedoes = source["torpedoes"];
	    }
	}
	export class SurvivedRate {
	    all: number;
	    win: number;
	    lose: number;
	
	    static createFrom(source: any = {}) {
	        return new SurvivedRate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.all = source["all"];
	        this.win = source["win"];
	        this.lose = source["lose"];
	    }
	}
	export class MaxDamage {
	    shipID: number;
	    shipName: string;
	    shipTier: number;
	    value: number;
	
	    static createFrom(source: any = {}) {
	        return new MaxDamage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.shipID = source["shipID"];
	        this.shipName = source["shipName"];
	        this.shipTier = source["shipTier"];
	        this.value = source["value"];
	    }
	}
	export class ShipStats {
	    battles: number;
	    damage: RatingValue;
	    maxDamage: MaxDamage;
	    winRate: RatingValue;
	    survivedRate: SurvivedRate;
	    kdRate: number;
	    kill: number;
	    exp: number;
	    pr: RatingValue;
	    hitRate: HitRate;
	    planesKilled: number;
	    platoonRate: number;
	    efficiencyBadge: string;
	
	    static createFrom(source: any = {}) {
	        return new ShipStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.battles = source["battles"];
	        this.damage = this.convertValues(source["damage"], RatingValue);
	        this.maxDamage = this.convertValues(source["maxDamage"], MaxDamage);
	        this.winRate = this.convertValues(source["winRate"], RatingValue);
	        this.survivedRate = this.convertValues(source["survivedRate"], SurvivedRate);
	        this.kdRate = source["kdRate"];
	        this.kill = source["kill"];
	        this.exp = source["exp"];
	        this.pr = this.convertValues(source["pr"], RatingValue);
	        this.hitRate = this.convertValues(source["hitRate"], HitRate);
	        this.planesKilled = source["planesKilled"];
	        this.platoonRate = source["platoonRate"];
	        this.efficiencyBadge = source["efficiencyBadge"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PlayerStats {
	    ship: ShipStats;
	    overall: OverallStats;
	
	    static createFrom(source: any = {}) {
	        return new PlayerStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ship = this.convertValues(source["ship"], ShipStats);
	        this.overall = this.convertValues(source["overall"], OverallStats);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class RatingValue {
	    value: number;
	    rating: string;
	
	    static createFrom(source: any = {}) {
	        return new RatingValue(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.value = source["value"];
	        this.rating = source["rating"];
	    }
	}
	export class ServerAverage {
	    damage: number;
	    frags: number;
	    winRate: number;
	
	    static createFrom(source: any = {}) {
	        return new ServerAverage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.damage = source["damage"];
	        this.frags = source["frags"];
	        this.winRate = source["winRate"];
	    }
	}
	export class Warship {
	    id: number;
	    name: string;
	    tier: number;
	    type: string;
	    nation: string;
	    isPremium: boolean;
	    serverAverage?: ServerAverage;
	    damageRatings?: RatingValue[];
	
	    static createFrom(source: any = {}) {
	        return new Warship(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.tier = source["tier"];
	        this.type = source["type"];
	        this.nation = source["nation"];
	        this.isPremium = source["isPremium"];
	        this.serverAverage = this.convertValues(source["serverAverage"], ServerAverage);
	        this.damageRatings = this.convertValues(source["damageRatings"], RatingValue);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Clan {
	    id: number;
	    tag: string;
	    colorCode: string;
	    language: string;
	
	    static createFrom(source: any = {}) {
	        return new Clan(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.tag = source["tag"];
	        this.colorCode = source["colorCode"];
	        this.language = source["language"];
	    }
	}
	export class PlayerInfo {
	    id: number;
	    name: string;
	    clan: Clan;
	    isHidden: boolean;
	
	    static createFrom(source: any = {}) {
	        return new PlayerInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.clan = this.convertValues(source["clan"], Clan);
	        this.isHidden = source["isHidden"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Player {
	    playerInfo: PlayerInfo;
	    warship: Warship;
	    pvpSolo: PlayerStats;
	    pvpAll: PlayerStats;
	    rankSolo: PlayerStats;
	
	    static createFrom(source: any = {}) {
	        return new Player(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.playerInfo = this.convertValues(source["playerInfo"], PlayerInfo);
	        this.warship = this.convertValues(source["warship"], Warship);
	        this.pvpSolo = this.convertValues(source["pvpSolo"], PlayerStats);
	        this.pvpAll = this.convertValues(source["pvpAll"], PlayerStats);
	        this.rankSolo = this.convertValues(source["rankSolo"], PlayerStats);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Team {
	    players: Player[];
	    pvpSolo: TeamStats;
	    pvpAll: TeamStats;
	    rankSolo: TeamStats;
	
	    static createFrom(source: any = {}) {
	        return new Team(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.players = this.convertValues(source["players"], Player);
	        this.pvpSolo = this.convertValues(source["pvpSolo"], TeamStats);
	        this.pvpAll = this.convertValues(source["pvpAll"], TeamStats);
	        this.rankSolo = this.convertValues(source["rankSolo"], TeamStats);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class BattleMetadata {
	    unixtime: number;
	    arena: string;
	    type: string;
	
	    static createFrom(source: any = {}) {
	        return new BattleMetadata(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.unixtime = source["unixtime"];
	        this.arena = source["arena"];
	        this.type = source["type"];
	    }
	}
	export class Battle {
	    metadata: BattleMetadata;
	    teams: Team[];
	
	    static createFrom(source: any = {}) {
	        return new Battle(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.metadata = this.convertValues(source["metadata"], BattleMetadata);
	        this.teams = this.convertValues(source["teams"], Team);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class DetailStatsColumnPref {
	    isShowShip: boolean;
	    isShowOverall: boolean;
	    digit: number;
	
	    static createFrom(source: any = {}) {
	        return new DetailStatsColumnPref(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.isShowShip = source["isShowShip"];
	        this.isShowOverall = source["isShowOverall"];
	        this.digit = source["digit"];
	    }
	}
	export class StatsColumnPref {
	    battles: DetailStatsColumnPref;
	    damage: DetailStatsColumnPref;
	    maxDamage: DetailStatsColumnPref;
	    winRate: DetailStatsColumnPref;
	    survivedRate: DetailStatsColumnPref;
	    kdRate: DetailStatsColumnPref;
	    kill: DetailStatsColumnPref;
	    exp: DetailStatsColumnPref;
	    pr: DetailStatsColumnPref;
	    hitRate: DetailStatsColumnPref;
	    planesKilled: DetailStatsColumnPref;
	    platoonRate: DetailStatsColumnPref;
	    efficiencyBadge: DetailStatsColumnPref;
	    threatLevel: DetailStatsColumnPref;
	    avgTier: DetailStatsColumnPref;
	    usingShipTypeRate: DetailStatsColumnPref;
	    usingTierRate: DetailStatsColumnPref;
	
	    static createFrom(source: any = {}) {
	        return new StatsColumnPref(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.battles = this.convertValues(source["battles"], DetailStatsColumnPref);
	        this.damage = this.convertValues(source["damage"], DetailStatsColumnPref);
	        this.maxDamage = this.convertValues(source["maxDamage"], DetailStatsColumnPref);
	        this.winRate = this.convertValues(source["winRate"], DetailStatsColumnPref);
	        this.survivedRate = this.convertValues(source["survivedRate"], DetailStatsColumnPref);
	        this.kdRate = this.convertValues(source["kdRate"], DetailStatsColumnPref);
	        this.kill = this.convertValues(source["kill"], DetailStatsColumnPref);
	        this.exp = this.convertValues(source["exp"], DetailStatsColumnPref);
	        this.pr = this.convertValues(source["pr"], DetailStatsColumnPref);
	        this.hitRate = this.convertValues(source["hitRate"], DetailStatsColumnPref);
	        this.planesKilled = this.convertValues(source["planesKilled"], DetailStatsColumnPref);
	        this.platoonRate = this.convertValues(source["platoonRate"], DetailStatsColumnPref);
	        this.efficiencyBadge = this.convertValues(source["efficiencyBadge"], DetailStatsColumnPref);
	        this.threatLevel = this.convertValues(source["threatLevel"], DetailStatsColumnPref);
	        this.avgTier = this.convertValues(source["avgTier"], DetailStatsColumnPref);
	        this.usingShipTypeRate = this.convertValues(source["usingShipTypeRate"], DetailStatsColumnPref);
	        this.usingTierRate = this.convertValues(source["usingTierRate"], DetailStatsColumnPref);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ShipColumnPref {
	    enableNationFlag: boolean;
	    isColored: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ShipColumnPref(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enableNationFlag = source["enableNationFlag"];
	        this.isColored = source["isColored"];
	    }
	}
	export class PlayerColumnPref {
	    enableNationFlag: boolean;
	    colorPattern: string;
	
	    static createFrom(source: any = {}) {
	        return new PlayerColumnPref(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enableNationFlag = source["enableNationFlag"];
	        this.colorPattern = source["colorPattern"];
	    }
	}
	export class ColumnPref {
	    player: PlayerColumnPref;
	    ship: ShipColumnPref;
	    stats: StatsColumnPref;
	
	    static createFrom(source: any = {}) {
	        return new ColumnPref(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.player = this.convertValues(source["player"], PlayerColumnPref);
	        this.ship = this.convertValues(source["ship"], ShipColumnPref);
	        this.stats = this.convertValues(source["stats"], StatsColumnPref);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	export class NewVersion {
	    version: string;
	    downloadURL: string;
	
	    static createFrom(source: any = {}) {
	        return new NewVersion(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.downloadURL = source["downloadURL"];
	    }
	}
	
	
	
	
	
	export class Pref {
	    version: number;
	    installPath: string;
	    zoomRate: number;
	    statsExtra: string;
	    isSendReport: boolean;
	    column: ColumnPref;
	
	    static createFrom(source: any = {}) {
	        return new Pref(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.installPath = source["installPath"];
	        this.zoomRate = source["zoomRate"];
	        this.statsExtra = source["statsExtra"];
	        this.isSendReport = source["isSendReport"];
	        this.column = this.convertValues(source["column"], ColumnPref);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	
	
	
	
	
	
	
	export class Vehicle {
	    shipId: number;
	    relation: number;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new Vehicle(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.shipId = source["shipId"];
	        this.relation = source["relation"];
	        this.name = source["name"];
	    }
	}
	export class TempArenaInfo {
	    vehicles: Vehicle[];
	    dateTime: string;
	    mapId: number;
	    matchGroup: string;
	    playerName: string;
	
	    static createFrom(source: any = {}) {
	        return new TempArenaInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.vehicles = this.convertValues(source["vehicles"], Vehicle);
	        this.dateTime = source["dateTime"];
	        this.mapId = source["mapId"];
	        this.matchGroup = source["matchGroup"];
	        this.playerName = source["playerName"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	

}

