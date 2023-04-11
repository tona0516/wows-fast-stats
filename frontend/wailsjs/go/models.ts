export namespace vo {
	
	export class PlayerPlayerStats {
	    battles: number;
	    avg_damage: number;
	    avg_exp: number;
	    win_rate: number;
	    kd_rate: number;
	    avg_tier: number;
	
	    static createFrom(source: any = {}) {
	        return new PlayerPlayerStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.battles = source["battles"];
	        this.avg_damage = source["avg_damage"];
	        this.avg_exp = source["avg_exp"];
	        this.win_rate = source["win_rate"];
	        this.kd_rate = source["kd_rate"];
	        this.avg_tier = source["avg_tier"];
	    }
	}
	export class PlayerPlayerInfo {
	    name: string;
	    clan: string;
	    is_hidden: boolean;
	    stats_url: string;
	
	    static createFrom(source: any = {}) {
	        return new PlayerPlayerInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.clan = source["clan"];
	        this.is_hidden = source["is_hidden"];
	        this.stats_url = source["stats_url"];
	    }
	}
	export class PlayerShipStats {
	    battles: number;
	    avg_damage: number;
	    avg_exp: number;
	    win_rate: number;
	    kd_rate: number;
	    combat_power: number;
	    personal_rating: number;
	
	    static createFrom(source: any = {}) {
	        return new PlayerShipStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.battles = source["battles"];
	        this.avg_damage = source["avg_damage"];
	        this.avg_exp = source["avg_exp"];
	        this.win_rate = source["win_rate"];
	        this.kd_rate = source["kd_rate"];
	        this.combat_power = source["combat_power"];
	        this.personal_rating = source["personal_rating"];
	    }
	}
	export class PlayerShipInfo {
	    name: string;
	    nation: string;
	    tier: number;
	    type: string;
	    stats_url: string;
	
	    static createFrom(source: any = {}) {
	        return new PlayerShipInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.nation = source["nation"];
	        this.tier = source["tier"];
	        this.type = source["type"];
	        this.stats_url = source["stats_url"];
	    }
	}
	export class Player {
	    player_ship_info: PlayerShipInfo;
	    player_ship_stats: PlayerShipStats;
	    player_player_info: PlayerPlayerInfo;
	    player_player_stats: PlayerPlayerStats;
	
	    static createFrom(source: any = {}) {
	        return new Player(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.player_ship_info = this.convertValues(source["player_ship_info"], PlayerShipInfo);
	        this.player_ship_stats = this.convertValues(source["player_ship_stats"], PlayerShipStats);
	        this.player_player_info = this.convertValues(source["player_player_info"], PlayerPlayerInfo);
	        this.player_player_stats = this.convertValues(source["player_player_stats"], PlayerPlayerStats);
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
	    friends: Player[];
	    enemies: Player[];
	
	    static createFrom(source: any = {}) {
	        return new Team(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.friends = this.convertValues(source["friends"], Player);
	        this.enemies = this.convertValues(source["enemies"], Player);
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
	
	    static createFrom(source: any = {}) {
	        return new UserConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.install_path = source["install_path"];
	        this.appid = source["appid"];
	    }
	}

}

