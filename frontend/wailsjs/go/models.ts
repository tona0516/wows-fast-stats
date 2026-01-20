export namespace core {
	
	export class TeamThreatLevel {
	    average: number;
	    dissociation_degree: number;
	    accuracy: number;
	
	    static createFrom(source: any = {}) {
	        return new TeamThreatLevel(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.average = source["average"];
	        this.dissociation_degree = source["dissociation_degree"];
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
	    team_average_stats: TeamAverageStats;
	    team_threat_level: TeamThreatLevel;
	
	    static createFrom(source: any = {}) {
	        return new TeamStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.team_average_stats = this.convertValues(source["team_average_stats"], TeamAverageStats);
	        this.team_threat_level = this.convertValues(source["team_threat_level"], TeamThreatLevel);
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
	    ss: number;
	    dd: number;
	    cl: number;
	    bb: number;
	    cv: number;
	
	    static createFrom(source: any = {}) {
	        return new ShipTypeGroup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ss = source["ss"];
	        this.dd = source["dd"];
	        this.cl = source["cl"];
	        this.bb = source["bb"];
	        this.cv = source["cv"];
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
	    max_damage: MaxDamage;
	    win_rate: RatingValue;
	    survived_rate: SurvivedRate;
	    kd_rate: number;
	    kill: number;
	    exp: number;
	    pr: RatingValue;
	    threat_level: ThreatLevel;
	    avg_tier: number;
	    using_ship_type_rate: ShipTypeGroup;
	    using_tier_rate: TierGroup;
	    platoon_rate: number;
	    efficiency_badge: EfficiencyBadgeGroup;
	
	    static createFrom(source: any = {}) {
	        return new OverallStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.battles = source["battles"];
	        this.damage = this.convertValues(source["damage"], RatingValue);
	        this.max_damage = this.convertValues(source["max_damage"], MaxDamage);
	        this.win_rate = this.convertValues(source["win_rate"], RatingValue);
	        this.survived_rate = this.convertValues(source["survived_rate"], SurvivedRate);
	        this.kd_rate = source["kd_rate"];
	        this.kill = source["kill"];
	        this.exp = source["exp"];
	        this.pr = this.convertValues(source["pr"], RatingValue);
	        this.threat_level = this.convertValues(source["threat_level"], ThreatLevel);
	        this.avg_tier = source["avg_tier"];
	        this.using_ship_type_rate = this.convertValues(source["using_ship_type_rate"], ShipTypeGroup);
	        this.using_tier_rate = this.convertValues(source["using_tier_rate"], TierGroup);
	        this.platoon_rate = source["platoon_rate"];
	        this.efficiency_badge = this.convertValues(source["efficiency_badge"], EfficiencyBadgeGroup);
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
	    main_battery: number;
	    torpedoes: number;
	
	    static createFrom(source: any = {}) {
	        return new HitRate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.main_battery = source["main_battery"];
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
	    ship_id: number;
	    ship_name: string;
	    ship_tier: number;
	    value: number;
	
	    static createFrom(source: any = {}) {
	        return new MaxDamage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ship_id = source["ship_id"];
	        this.ship_name = source["ship_name"];
	        this.ship_tier = source["ship_tier"];
	        this.value = source["value"];
	    }
	}
	export class ShipStats {
	    battles: number;
	    damage: RatingValue;
	    max_damage: MaxDamage;
	    win_rate: RatingValue;
	    survived_rate: SurvivedRate;
	    kd_rate: number;
	    kill: number;
	    exp: number;
	    pr: RatingValue;
	    hit_rate: HitRate;
	    planes_killed: number;
	    platoon_rate: number;
	    efficiency_badge: string;
	
	    static createFrom(source: any = {}) {
	        return new ShipStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.battles = source["battles"];
	        this.damage = this.convertValues(source["damage"], RatingValue);
	        this.max_damage = this.convertValues(source["max_damage"], MaxDamage);
	        this.win_rate = this.convertValues(source["win_rate"], RatingValue);
	        this.survived_rate = this.convertValues(source["survived_rate"], SurvivedRate);
	        this.kd_rate = source["kd_rate"];
	        this.kill = source["kill"];
	        this.exp = source["exp"];
	        this.pr = this.convertValues(source["pr"], RatingValue);
	        this.hit_rate = this.convertValues(source["hit_rate"], HitRate);
	        this.planes_killed = source["planes_killed"];
	        this.platoon_rate = source["platoon_rate"];
	        this.efficiency_badge = source["efficiency_badge"];
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
	    hex_color: string;
	    language: string;
	
	    static createFrom(source: any = {}) {
	        return new Clan(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.tag = source["tag"];
	        this.hex_color = source["hex_color"];
	        this.language = source["language"];
	    }
	}
	export class PlayerInfo {
	    id: number;
	    name: string;
	    clan: Clan;
	    is_hidden: boolean;
	
	    static createFrom(source: any = {}) {
	        return new PlayerInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.clan = this.convertValues(source["clan"], Clan);
	        this.is_hidden = source["is_hidden"];
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
	    player_info: PlayerInfo;
	    warship: Warship;
	    pvp_solo: PlayerStats;
	    pvp_all: PlayerStats;
	    rank_solo: PlayerStats;
	
	    static createFrom(source: any = {}) {
	        return new Player(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.player_info = this.convertValues(source["player_info"], PlayerInfo);
	        this.warship = this.convertValues(source["warship"], Warship);
	        this.pvp_solo = this.convertValues(source["pvp_solo"], PlayerStats);
	        this.pvp_all = this.convertValues(source["pvp_all"], PlayerStats);
	        this.rank_solo = this.convertValues(source["rank_solo"], PlayerStats);
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
	    pvp_solo: TeamStats;
	    pvp_all: TeamStats;
	    rank_solo: TeamStats;
	
	    static createFrom(source: any = {}) {
	        return new Team(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.players = this.convertValues(source["players"], Player);
	        this.pvp_solo = this.convertValues(source["pvp_solo"], TeamStats);
	        this.pvp_all = this.convertValues(source["pvp_all"], TeamStats);
	        this.rank_solo = this.convertValues(source["rank_solo"], TeamStats);
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
	
	
	export class DetailStatsColumnConfig {
	    is_show_ship: boolean;
	    is_show_overall: boolean;
	    digit: number;
	
	    static createFrom(source: any = {}) {
	        return new DetailStatsColumnConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.is_show_ship = source["is_show_ship"];
	        this.is_show_overall = source["is_show_overall"];
	        this.digit = source["digit"];
	    }
	}
	export class StatsColumnConfig {
	    battles: DetailStatsColumnConfig;
	    damage: DetailStatsColumnConfig;
	    max_damage: DetailStatsColumnConfig;
	    win_rate: DetailStatsColumnConfig;
	    survived_rate: DetailStatsColumnConfig;
	    kd_rate: DetailStatsColumnConfig;
	    kill: DetailStatsColumnConfig;
	    exp: DetailStatsColumnConfig;
	    pr: DetailStatsColumnConfig;
	    hit_rate: DetailStatsColumnConfig;
	    planes_killed: DetailStatsColumnConfig;
	    platoon_rate: DetailStatsColumnConfig;
	    efficiency_badge: DetailStatsColumnConfig;
	    threat_level: DetailStatsColumnConfig;
	    avg_tier: DetailStatsColumnConfig;
	    using_ship_type_rate: DetailStatsColumnConfig;
	    using_tier_rate: DetailStatsColumnConfig;
	
	    static createFrom(source: any = {}) {
	        return new StatsColumnConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.battles = this.convertValues(source["battles"], DetailStatsColumnConfig);
	        this.damage = this.convertValues(source["damage"], DetailStatsColumnConfig);
	        this.max_damage = this.convertValues(source["max_damage"], DetailStatsColumnConfig);
	        this.win_rate = this.convertValues(source["win_rate"], DetailStatsColumnConfig);
	        this.survived_rate = this.convertValues(source["survived_rate"], DetailStatsColumnConfig);
	        this.kd_rate = this.convertValues(source["kd_rate"], DetailStatsColumnConfig);
	        this.kill = this.convertValues(source["kill"], DetailStatsColumnConfig);
	        this.exp = this.convertValues(source["exp"], DetailStatsColumnConfig);
	        this.pr = this.convertValues(source["pr"], DetailStatsColumnConfig);
	        this.hit_rate = this.convertValues(source["hit_rate"], DetailStatsColumnConfig);
	        this.planes_killed = this.convertValues(source["planes_killed"], DetailStatsColumnConfig);
	        this.platoon_rate = this.convertValues(source["platoon_rate"], DetailStatsColumnConfig);
	        this.efficiency_badge = this.convertValues(source["efficiency_badge"], DetailStatsColumnConfig);
	        this.threat_level = this.convertValues(source["threat_level"], DetailStatsColumnConfig);
	        this.avg_tier = this.convertValues(source["avg_tier"], DetailStatsColumnConfig);
	        this.using_ship_type_rate = this.convertValues(source["using_ship_type_rate"], DetailStatsColumnConfig);
	        this.using_tier_rate = this.convertValues(source["using_tier_rate"], DetailStatsColumnConfig);
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
	export class ShipColumnConfig {
	    enable_nation_flag: boolean;
	    is_colored: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ShipColumnConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enable_nation_flag = source["enable_nation_flag"];
	        this.is_colored = source["is_colored"];
	    }
	}
	export class PlayerColumnConfig {
	    enable_nation_flag: boolean;
	    color_pattern: string;
	
	    static createFrom(source: any = {}) {
	        return new PlayerColumnConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enable_nation_flag = source["enable_nation_flag"];
	        this.color_pattern = source["color_pattern"];
	    }
	}
	export class ColumnConfig {
	    player: PlayerColumnConfig;
	    ship: ShipColumnConfig;
	    stats: StatsColumnConfig;
	
	    static createFrom(source: any = {}) {
	        return new ColumnConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.player = this.convertValues(source["player"], PlayerColumnConfig);
	        this.ship = this.convertValues(source["ship"], ShipColumnConfig);
	        this.stats = this.convertValues(source["stats"], StatsColumnConfig);
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
	    url: string;
	
	    static createFrom(source: any = {}) {
	        return new NewVersion(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.url = source["url"];
	    }
	}
	
	
	
	
	
	export class Pref {
	    version: number;
	    install_path: string;
	    zoom_rate: number;
	    stats_extra: string;
	    is_send_report: boolean;
	    column: ColumnConfig;
	
	    static createFrom(source: any = {}) {
	        return new Pref(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.install_path = source["install_path"];
	        this.zoom_rate = source["zoom_rate"];
	        this.stats_extra = source["stats_extra"];
	        this.is_send_report = source["is_send_report"];
	        this.column = this.convertValues(source["column"], ColumnConfig);
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

