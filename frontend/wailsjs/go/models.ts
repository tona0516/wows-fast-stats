export namespace data {
	
	export class AlertPlayer {
	    account_id: number;
	    name: string;
	    pattern: string;
	    message: string;
	    created_at: number;
	
	    static createFrom(source: any = {}) {
	        return new AlertPlayer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.account_id = source["account_id"];
	        this.name = source["name"];
	        this.pattern = source["pattern"];
	        this.message = source["message"];
	        this.created_at = source["created_at"];
	    }
	}
	export class PlayerColumnSetting {
	    enable_nation_flag: boolean;
	    color_pattern: string;
	
	    static createFrom(source: any = {}) {
	        return new PlayerColumnSetting(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enable_nation_flag = source["enable_nation_flag"];
	        this.color_pattern = source["color_pattern"];
	    }
	}
	export class ShipColumnSetting {
	    enable_nation_flag: boolean;
	    is_colored: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ShipColumnSetting(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enable_nation_flag = source["enable_nation_flag"];
	        this.is_colored = source["is_colored"];
	    }
	}
	export class BasicColumnSetting {
	    version: number;
	    ship: ShipColumnSetting;
	    player: PlayerColumnSetting;
	
	    static createFrom(source: any = {}) {
	        return new BasicColumnSetting(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.ship = this.convertValues(source["ship"], ShipColumnSetting);
	        this.player = this.convertValues(source["player"], PlayerColumnSetting);
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
	    raw: number;
	    modified: number;
	
	    static createFrom(source: any = {}) {
	        return new ThreatLevel(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
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
	
	    static createFrom(source: any = {}) {
	        return new Team(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.players = this.convertValues(source["players"], Player);
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
	
	
	
	export class GHLatestRelease {
	    tag_name: string;
	    html_url: string;
	    updatable: boolean;
	
	    static createFrom(source: any = {}) {
	        return new GHLatestRelease(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tag_name = source["tag_name"];
	        this.html_url = source["html_url"];
	        this.updatable = source["updatable"];
	    }
	}
	
	
	export class OptionalSetting {
	    version: number;
	    zoom_rate: number;
	    stats_extra: string;
	    is_send_report: boolean;
	
	    static createFrom(source: any = {}) {
	        return new OptionalSetting(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.zoom_rate = source["zoom_rate"];
	        this.stats_extra = source["stats_extra"];
	        this.is_send_report = source["is_send_report"];
	    }
	}
	
	
	
	
	
	
	export class RequiredSetting {
	    version: number;
	    install_path: string;
	
	    static createFrom(source: any = {}) {
	        return new RequiredSetting(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.install_path = source["install_path"];
	    }
	}
	
	
	
	
	export class StatsColumnSetting {
	    is_show_ship: boolean;
	    is_show_overall: boolean;
	    digit: number;
	
	    static createFrom(source: any = {}) {
	        return new StatsColumnSetting(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.is_show_ship = source["is_show_ship"];
	        this.is_show_overall = source["is_show_overall"];
	        this.digit = source["digit"];
	    }
	}
	export class StatsColumnSettings {
	    version: number;
	    battles: StatsColumnSetting;
	    damage: StatsColumnSetting;
	    max_damage: StatsColumnSetting;
	    win_rate: StatsColumnSetting;
	    survived_rate: StatsColumnSetting;
	    kd_rate: StatsColumnSetting;
	    kill: StatsColumnSetting;
	    exp: StatsColumnSetting;
	    pr: StatsColumnSetting;
	    hit_rate: StatsColumnSetting;
	    planes_killed: StatsColumnSetting;
	    platoon_rate: StatsColumnSetting;
	    efficiency_badge: StatsColumnSetting;
	    threat_level: StatsColumnSetting;
	    avg_tier: StatsColumnSetting;
	    using_ship_type_rate: StatsColumnSetting;
	    using_tier_rate: StatsColumnSetting;
	
	    static createFrom(source: any = {}) {
	        return new StatsColumnSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.battles = this.convertValues(source["battles"], StatsColumnSetting);
	        this.damage = this.convertValues(source["damage"], StatsColumnSetting);
	        this.max_damage = this.convertValues(source["max_damage"], StatsColumnSetting);
	        this.win_rate = this.convertValues(source["win_rate"], StatsColumnSetting);
	        this.survived_rate = this.convertValues(source["survived_rate"], StatsColumnSetting);
	        this.kd_rate = this.convertValues(source["kd_rate"], StatsColumnSetting);
	        this.kill = this.convertValues(source["kill"], StatsColumnSetting);
	        this.exp = this.convertValues(source["exp"], StatsColumnSetting);
	        this.pr = this.convertValues(source["pr"], StatsColumnSetting);
	        this.hit_rate = this.convertValues(source["hit_rate"], StatsColumnSetting);
	        this.planes_killed = this.convertValues(source["planes_killed"], StatsColumnSetting);
	        this.platoon_rate = this.convertValues(source["platoon_rate"], StatsColumnSetting);
	        this.efficiency_badge = this.convertValues(source["efficiency_badge"], StatsColumnSetting);
	        this.threat_level = this.convertValues(source["threat_level"], StatsColumnSetting);
	        this.avg_tier = this.convertValues(source["avg_tier"], StatsColumnSetting);
	        this.using_ship_type_rate = this.convertValues(source["using_ship_type_rate"], StatsColumnSetting);
	        this.using_tier_rate = this.convertValues(source["using_tier_rate"], StatsColumnSetting);
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

