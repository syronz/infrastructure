import { Component, OnInit } from '@angular/core';
import { trigger, state, style, animate, transition } from '@angular/animations';
import { DictionaryService } from '../../services/dictionary.service';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-sidebar',
	templateUrl: './sidebar.component.html',
	styleUrls: ['./sidebar.component.css'],
	animations: [
		trigger('menuState', [
			state('open', style({
				height: '*'
			})),
			state('close', style({
				height: '0',
				padding: '0'
			})),
			transition('open <=> close', animate('150ms ease-in')),
		])		

	]
})
export class SidebarComponent implements OnInit {
	sidebarItems = {
		manage: 'close',
		reports: 'close',
	}
	userAcl: any;

	constructor(
		public dict: DictionaryService,
		private authService: AuthService,
	) { }

	ngOnInit() {
		this.dict.loadWords();

		this.userAcl = this.authService.getUserJWT()['acls'];
	}

	open(str) {
		console.info(str);
	}

	toggleState(menuPart) {
		this.sidebarItems[menuPart] = this.sidebarItems[menuPart] === 'open' ? 'close' : 'open';
	}

}
