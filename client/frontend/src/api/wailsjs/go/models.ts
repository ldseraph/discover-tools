export namespace main {
	
	export class IP4Address {
	    address: string;
	    prefix: number;
	    gateway: string;
	
	    static createFrom(source: any = {}) {
	        return new IP4Address(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.address = source["address"];
	        this.prefix = source["prefix"];
	        this.gateway = source["gateway"];
	    }
	}
	export class Node {
	    uuid: string;
	    ip4: IP4Address[];
	    dhcp: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Node(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uuid = source["uuid"];
	        this.ip4 = this.convertValues(source["ip4"], IP4Address);
	        this.dhcp = source["dhcp"];
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

