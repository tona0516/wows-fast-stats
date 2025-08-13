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
	    damage: number;
	    max_damage: MaxDamage;
	    win_rate: number;
	    win_survived_rate: number;
	    lose_survived_rate: number;
	    kd_rate: number;
	    kill: number;
	    exp: number;
	    pr: number;
	    threat_level: ThreatLevel;
	    avg_tier: number;
	    using_ship_type_rate: ShipTypeGroup;
	    using_tier_rate: TierGroup;
	    platoon_rate: number;
	
	    static createFrom(source: any = {}) {
	        return new OverallStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.battles = source["battles"];
	        this.damage = source["damage"];
	        this.max_damage = this.convertValues(source["max_damage"], MaxDamage);
	        this.win_rate = source["win_rate"];
	        this.win_survived_rate = source["win_survived_rate"];
	        this.lose_survived_rate = source["lose_survived_rate"];
	        this.kd_rate = source["kd_rate"];
	        this.kill = source["kill"];
	        this.exp = source["exp"];
	        this.pr = source["pr"];
	        this.threat_level = this.convertValues(source["threat_level"], ThreatLevel);
	        this.avg_tier = source["avg_tier"];
	        this.using_ship_type_rate = this.convertValues(source["using_ship_type_rate"], ShipTypeGroup);
	        this.using_tier_rate = this.convertValues(source["using_tier_rate"], TierGroup);
	        this.platoon_rate = source["platoon_rate"];
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
	    damage: number;
	    max_damage: MaxDamage;
	    win_rate: number;
	    win_survived_rate: number;
	    lose_survived_rate: number;
	    kd_rate: number;
	    kill: number;
	    exp: number;
	    pr: number;
	    main_battery_hit_rate: number;
	    torpedoes_hit_rate: number;
	    planes_killed: number;
	    platoon_rate: number;
	
	    static createFrom(source: any = {}) {
	        return new ShipStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.battles = source["battles"];
	        this.damage = source["damage"];
	        this.max_damage = this.convertValues(source["max_damage"], MaxDamage);
	        this.win_rate = source["win_rate"];
	        this.win_survived_rate = source["win_survived_rate"];
	        this.lose_survived_rate = source["lose_survived_rate"];
	        this.kd_rate = source["kd_rate"];
	        this.kill = source["kill"];
	        this.exp = source["exp"];
	        this.pr = source["pr"];
	        this.main_battery_hit_rate = source["main_battery_hit_rate"];
	        this.torpedoes_hit_rate = source["torpedoes_hit_rate"];
	        this.planes_killed = source["planes_killed"];
	        this.platoon_rate = source["platoon_rate"];
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
	export class ShipInfo {
	    id: number;
	    name: string;
	    nation: string;
	    tier: number;
	    type: string;
	    is_premium: boolean;
	    avg_damage: number;
	
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
	export class Meta {
	    unixtime: number;
	    arena: string;
	    type: string;
	    own_ship: string;
	
	    static createFrom(source: any = {}) {
	        return new Meta(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.unixtime = source["unixtime"];
	        this.arena = source["arena"];
	        this.type = source["type"];
	        this.own_ship = source["own_ship"];
	    }
	}
	export class Battle {
	    meta: Meta;
	    teams: Team[];
	
	    static createFrom(source: any = {}) {
	        return new Battle(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.meta = this.convertValues(source["meta"], Meta);
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

