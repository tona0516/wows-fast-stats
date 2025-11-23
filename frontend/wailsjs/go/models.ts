export namespace data {
	
	export class TeamAverageStats {
	    ship_avg_pr: number;
	    ship_avg_damage: number;
	    ship_win_rate: number;
	    overall_avg_pr: number;
	    overall_avg_damage: number;
	    overall_win_rate: number;
	
	    static createFrom(source: any = {}) {
	        return new TeamAverageStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ship_avg_pr = source["ship_avg_pr"];
	        this.ship_avg_damage = source["ship_avg_damage"];
	        this.ship_win_rate = source["ship_win_rate"];
	        this.overall_avg_pr = source["overall_avg_pr"];
	        this.overall_avg_damage = source["overall_avg_damage"];
	        this.overall_win_rate = source["overall_win_rate"];
	    }
	}
	export class TeamShipTypeStats {
	    cv: TeamAverageStats;
	    bb: TeamAverageStats;
	    cl: TeamAverageStats;
	    dd: TeamAverageStats;
	    ss: TeamAverageStats;
	
	    static createFrom(source: any = {}) {
	        return new TeamShipTypeStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cv = this.convertValues(source["cv"], TeamAverageStats);
	        this.bb = this.convertValues(source["bb"], TeamAverageStats);
	        this.cl = this.convertValues(source["cl"], TeamAverageStats);
	        this.dd = this.convertValues(source["dd"], TeamAverageStats);
	        this.ss = this.convertValues(source["ss"], TeamAverageStats);
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
	export class TeamStats {
	    team_threat_level: TeamThreatLevel;
	
	    static createFrom(source: any = {}) {
	        return new TeamStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
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
	export class ShipInfo {
	    id: number;
	    name: string;
	    nation: string;
	    tier: number;
	    type: string;
	    is_premium: boolean;
	    avg_damage: number;
	    damage_ratings: RatingValue[];
	
	    static createFrom(source: any = {}) {
	        return new ShipInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.nation = source["nation"];
	        this.tier = source["tier"];
	        this.type = source["type"];
	        this.is_premium = source["is_premium"];
	        this.avg_damage = source["avg_damage"];
	        this.damage_ratings = this.convertValues(source["damage_ratings"], RatingValue);
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
	    tag: string;
	    id: number;
	    hex_color: string;
	    language: string;
	
	    static createFrom(source: any = {}) {
	        return new Clan(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tag = source["tag"];
	        this.id = source["id"];
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
	    ship_info: ShipInfo;
	    pvp_solo: PlayerStats;
	    pvp_all: PlayerStats;
	    rank_solo: PlayerStats;
	
	    static createFrom(source: any = {}) {
	        return new Player(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.player_info = this.convertValues(source["player_info"], PlayerInfo);
	        this.ship_info = this.convertValues(source["ship_info"], ShipInfo);
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
	    ship_type_stats: TeamShipTypeStats;
	
	    static createFrom(source: any = {}) {
	        return new Team(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.players = this.convertValues(source["players"], Player);
	        this.pvp_solo = this.convertValues(source["pvp_solo"], TeamStats);
	        this.pvp_all = this.convertValues(source["pvp_all"], TeamStats);
	        this.rank_solo = this.convertValues(source["rank_solo"], TeamStats);
	        this.ship_type_stats = this.convertValues(source["ship_type_stats"], TeamShipTypeStats);
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
	export class BattleMetaData {
	    unixtime: number;
	    arena: string;
	    type: string;
	
	    static createFrom(source: any = {}) {
	        return new BattleMetaData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.unixtime = source["unixtime"];
	        this.arena = source["arena"];
	        this.type = source["type"];
	    }
	}
	export class Battle {
	    metadata: BattleMetaData;
	    teams: Team[];
	
	    static createFrom(source: any = {}) {
	        return new Battle(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.metadata = this.convertValues(source["metadata"], BattleMetaData);
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
	    semver: string;
	    url: string;
	
	    static createFrom(source: any = {}) {
	        return new NewVersion(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.semver = source["semver"];
	        this.url = source["url"];
	    }
	}
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	export class UserConfig {
	    version: number;
	    install_path: string;
	    zoom_rate: number;
	    stats_extra: string;
	    is_send_report: boolean;
	    column: ColumnConfig;
	
	    static createFrom(source: any = {}) {
	        return new UserConfig(source);
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
	export class WGAccountListData {
	    nickname: string;
	    account_id: number;
	
	    static createFrom(source: any = {}) {
	        return new WGAccountListData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.nickname = source["nickname"];
	        this.account_id = source["account_id"];
	    }
	}

}

