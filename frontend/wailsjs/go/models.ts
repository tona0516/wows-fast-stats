export namespace vo {
	
	export class AlertPlayer {
	    account_id: number;
	    name: string;
	    pattern: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new AlertPlayer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.account_id = source["account_id"];
	        this.name = source["name"];
	        this.pattern = source["pattern"];
	        this.message = source["message"];
	    }
	}
	export class Basic {
	    is_in_avg: boolean;
	    player_name: boolean;
	    ship_info: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Basic(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.is_in_avg = source["is_in_avg"];
	        this.player_name = source["player_name"];
	        this.ship_info = source["ship_info"];
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
	export class OverallStats {
	    battles: number;
	    damage: number;
	    win_rate: number;
	    win_survived_rate: number;
	    lose_survived_rate: number;
	    kd_rate: number;
	    kill: number;
	    exp: number;
	    avg_tier: number;
	    using_ship_type_rate: ShipTypeGroup;
	    using_tier_rate: TierGroup;
	
	    static createFrom(source: any = {}) {
	        return new OverallStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.battles = source["battles"];
	        this.damage = source["damage"];
	        this.win_rate = source["win_rate"];
	        this.win_survived_rate = source["win_survived_rate"];
	        this.lose_survived_rate = source["lose_survived_rate"];
	        this.kd_rate = source["kd_rate"];
	        this.kill = source["kill"];
	        this.exp = source["exp"];
	        this.avg_tier = source["avg_tier"];
	        this.using_ship_type_rate = this.convertValues(source["using_ship_type_rate"], ShipTypeGroup);
	        this.using_tier_rate = this.convertValues(source["using_tier_rate"], TierGroup);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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
	export class ShipStats {
	    battles: number;
	    damage: number;
	    win_rate: number;
	    win_survived_rate: number;
	    lose_survived_rate: number;
	    kd_rate: number;
	    kill: number;
	    exp: number;
	    main_battery_hit_rate: number;
	    torpedoes_hit_rate: number;
	    planes_killed: number;
	    pr: number;
	
	    static createFrom(source: any = {}) {
	        return new ShipStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.battles = source["battles"];
	        this.damage = source["damage"];
	        this.win_rate = source["win_rate"];
	        this.win_survived_rate = source["win_survived_rate"];
	        this.lose_survived_rate = source["lose_survived_rate"];
	        this.kd_rate = source["kd_rate"];
	        this.kill = source["kill"];
	        this.exp = source["exp"];
	        this.main_battery_hit_rate = source["main_battery_hit_rate"];
	        this.torpedoes_hit_rate = source["torpedoes_hit_rate"];
	        this.planes_killed = source["planes_killed"];
	        this.pr = source["pr"];
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
		    if (a.slice) {
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
	        this.avg_damage = source["avg_damage"];
	    }
	}
	export class Clan {
	    tag: string;
	    id: number;
	
	    static createFrom(source: any = {}) {
	        return new Clan(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tag = source["tag"];
	        this.id = source["id"];
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
		    if (a.slice) {
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
	
	    static createFrom(source: any = {}) {
	        return new Player(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.player_info = this.convertValues(source["player_info"], PlayerInfo);
	        this.ship_info = this.convertValues(source["ship_info"], ShipInfo);
	        this.pvp_solo = this.convertValues(source["pvp_solo"], PlayerStats);
	        this.pvp_all = this.convertValues(source["pvp_all"], PlayerStats);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new Team(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.players = this.convertValues(source["players"], Player);
	        this.name = source["name"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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
	    date: string;
	    arena: string;
	    type: string;
	    own_ship: string;
	
	    static createFrom(source: any = {}) {
	        return new Meta(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date = source["date"];
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
		    if (a.slice) {
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
	
	export class Overall {
	    damage: boolean;
	    win_rate: boolean;
	    kd_rate: boolean;
	    kill: boolean;
	    exp: boolean;
	    battles: boolean;
	    survived_rate: boolean;
	    avg_tier: boolean;
	    using_ship_type_rate: boolean;
	    using_tier_rate: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Overall(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.damage = source["damage"];
	        this.win_rate = source["win_rate"];
	        this.kd_rate = source["kd_rate"];
	        this.kill = source["kill"];
	        this.exp = source["exp"];
	        this.battles = source["battles"];
	        this.survived_rate = source["survived_rate"];
	        this.avg_tier = source["avg_tier"];
	        this.using_ship_type_rate = source["using_ship_type_rate"];
	        this.using_tier_rate = source["using_tier_rate"];
	    }
	}
	export class Ship {
	    pr: boolean;
	    damage: boolean;
	    win_rate: boolean;
	    kd_rate: boolean;
	    kill: boolean;
	    exp: boolean;
	    battles: boolean;
	    survived_rate: boolean;
	    hit_rate: boolean;
	    planes_killed: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Ship(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pr = source["pr"];
	        this.damage = source["damage"];
	        this.win_rate = source["win_rate"];
	        this.kd_rate = source["kd_rate"];
	        this.kill = source["kill"];
	        this.exp = source["exp"];
	        this.battles = source["battles"];
	        this.survived_rate = source["survived_rate"];
	        this.hit_rate = source["hit_rate"];
	        this.planes_killed = source["planes_killed"];
	    }
	}
	export class Displays {
	    basic: Basic;
	    ship: Ship;
	    overall: Overall;
	
	    static createFrom(source: any = {}) {
	        return new Displays(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.basic = this.convertValues(source["basic"], Basic);
	        this.ship = this.convertValues(source["ship"], Ship);
	        this.overall = this.convertValues(source["overall"], Overall);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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
	
	
	
	
	
	
	
	
	
	
	
	
	export class UserConfig {
	    install_path: string;
	    appid: string;
	    font_size: string;
	    displays: Displays;
	    save_screenshot: boolean;
	    save_temp_arena_info: boolean;
	    stats_pattern: string;
	    send_report: boolean;
	
	    static createFrom(source: any = {}) {
	        return new UserConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.install_path = source["install_path"];
	        this.appid = source["appid"];
	        this.font_size = source["font_size"];
	        this.displays = this.convertValues(source["displays"], Displays);
	        this.save_screenshot = source["save_screenshot"];
	        this.save_temp_arena_info = source["save_temp_arena_info"];
	        this.stats_pattern = source["stats_pattern"];
	        this.send_report = source["send_report"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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
	export class Version {
	    semver: string;
	    revision: string;
	
	    static createFrom(source: any = {}) {
	        return new Version(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.semver = source["semver"];
	        this.revision = source["revision"];
	    }
	}
	export class WGError {
	    code: number;
	    message: string;
	    field: string;
	    value: string;
	
	    static createFrom(source: any = {}) {
	        return new WGError(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.message = source["message"];
	        this.field = source["field"];
	        this.value = source["value"];
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
	export class WGAccountList {
	    status: string;
	    data: WGAccountListData[];
	    error: WGError;
	
	    static createFrom(source: any = {}) {
	        return new WGAccountList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.data = this.convertValues(source["data"], WGAccountListData);
	        this.error = this.convertValues(source["error"], WGError);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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

