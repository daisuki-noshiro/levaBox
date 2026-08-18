export namespace metadata {
	
	export enum Source {
	    Bangumi = "bangumi",
	    VNDB = "vndb",
	}
	export class ImageCandidate {
	    Source: Source;
	    URL: string;
	    Thumbnail?: string;
	    Width?: number;
	    Height?: number;
	
	    static createFrom(source: any = {}) {
	        return new ImageCandidate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Source = source["Source"];
	        this.URL = source["URL"];
	        this.Thumbnail = source["Thumbnail"];
	        this.Width = source["Width"];
	        this.Height = source["Height"];
	    }
	}

}

export namespace model {
	
	export class Background {
	    Type: string;
	    Path: string;
	
	    static createFrom(source: any = {}) {
	        return new Background(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Type = source["Type"];
	        this.Path = source["Path"];
	    }
	}
	export class LaunchConfig {
	    ExecutablePath: string;
	    WorkingDirectory: string;
	
	    static createFrom(source: any = {}) {
	        return new LaunchConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ExecutablePath = source["ExecutablePath"];
	        this.WorkingDirectory = source["WorkingDirectory"];
	    }
	}
	export class Game {
	    ID: string;
	    Launch: LaunchConfig;
	    Title: string;
	    Company: string;
	    Year?: number;
	    Description: string;
	    CoverPath: string;
	    Background: Background;
	    BGMPath: string;
	    BGMEnabled: boolean;
	    Favorite: boolean;
	    Progress: string;
	    TotalPlaySeconds: number;
	    // Go type: time
	    LastPlayedAt?: any;
	    TagIDs: string[];
	
	    static createFrom(source: any = {}) {
	        return new Game(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Launch = this.convertValues(source["Launch"], LaunchConfig);
	        this.Title = source["Title"];
	        this.Company = source["Company"];
	        this.Year = source["Year"];
	        this.Description = source["Description"];
	        this.CoverPath = source["CoverPath"];
	        this.Background = this.convertValues(source["Background"], Background);
	        this.BGMPath = source["BGMPath"];
	        this.BGMEnabled = source["BGMEnabled"];
	        this.Favorite = source["Favorite"];
	        this.Progress = source["Progress"];
	        this.TotalPlaySeconds = source["TotalPlaySeconds"];
	        this.LastPlayedAt = this.convertValues(source["LastPlayedAt"], null);
	        this.TagIDs = source["TagIDs"];
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

export namespace service {
	
	export class ResolvedSource {
	    Source: metadata.Source;
	    ExternalID: string;
	
	    static createFrom(source: any = {}) {
	        return new ResolvedSource(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Source = source["Source"];
	        this.ExternalID = source["ExternalID"];
	    }
	}
	export class ImportDraft {
	    ExecutablePath: string;
	    WorkingDirectory: string;
	    SearchKeyword: string;
	    Title: string;
	    Company: string;
	    Year?: number;
	    Description?: string;
	    TagCandidates: string[];
	    CoverCandidates: metadata.ImageCandidate[];
	    BackgroundCandidates: metadata.ImageCandidate[];
	    Sources: ResolvedSource[];
	
	    static createFrom(source: any = {}) {
	        return new ImportDraft(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ExecutablePath = source["ExecutablePath"];
	        this.WorkingDirectory = source["WorkingDirectory"];
	        this.SearchKeyword = source["SearchKeyword"];
	        this.Title = source["Title"];
	        this.Company = source["Company"];
	        this.Year = source["Year"];
	        this.Description = source["Description"];
	        this.TagCandidates = source["TagCandidates"];
	        this.CoverCandidates = this.convertValues(source["CoverCandidates"], metadata.ImageCandidate);
	        this.BackgroundCandidates = this.convertValues(source["BackgroundCandidates"], metadata.ImageCandidate);
	        this.Sources = this.convertValues(source["Sources"], ResolvedSource);
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
	export class MetadataSourceIssue {
	    Source: metadata.Source;
	    Message: string;
	
	    static createFrom(source: any = {}) {
	        return new MetadataSourceIssue(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Source = source["Source"];
	        this.Message = source["Message"];
	    }
	}
	export class ImportMetadataResult {
	    Draft: ImportDraft;
	    Issues: MetadataSourceIssue[];
	
	    static createFrom(source: any = {}) {
	        return new ImportMetadataResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Draft = this.convertValues(source["Draft"], ImportDraft);
	        this.Issues = this.convertValues(source["Issues"], MetadataSourceIssue);
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
	
	
	export class SaveImportRequest {
	    ExecutablePath: string;
	    WorkingDirectory: string;
	    Title: string;
	    Company: string;
	    Year?: number;
	    Description?: string;
	    Tags: string[];
	    Cover?: metadata.ImageCandidate;
	    Background?: metadata.ImageCandidate;
	    Sources: ResolvedSource[];
	
	    static createFrom(source: any = {}) {
	        return new SaveImportRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ExecutablePath = source["ExecutablePath"];
	        this.WorkingDirectory = source["WorkingDirectory"];
	        this.Title = source["Title"];
	        this.Company = source["Company"];
	        this.Year = source["Year"];
	        this.Description = source["Description"];
	        this.Tags = source["Tags"];
	        this.Cover = this.convertValues(source["Cover"], metadata.ImageCandidate);
	        this.Background = this.convertValues(source["Background"], metadata.ImageCandidate);
	        this.Sources = this.convertValues(source["Sources"], ResolvedSource);
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
	export class StartImportResult {
	    status: string;
	    draft?: ImportDraft;
	    existingGameId?: string;
	
	    static createFrom(source: any = {}) {
	        return new StartImportResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.draft = this.convertValues(source["draft"], ImportDraft);
	        this.existingGameId = source["existingGameId"];
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

