import { Component, OnInit, EventEmitter, Output } from '@angular/core';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.css']
})
export class NavbarComponent implements OnInit {
	@Output() sideStatus = new EventEmitter<boolean>();
	currentSideStatus: boolean;
	userName: any;

	constructor(
		private authService: AuthService,
	) { 
		this.userName = authService.getUserJWT()
	}


  ngOnInit() {
  }

	toggleNav = function(){
		this.currentSideStatus = !this.currentSideStatus;
		this.sideStatus.emit(this.currentSideStatus);
	}

}
