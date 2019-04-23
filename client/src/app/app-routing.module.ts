import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { LoginComponent } from './components/login/login.component';
import { DashboardComponent } from './components/dashboard/dashboard.component';
import { PrivateFrameComponent } from './components/private-frame/private-frame.component';
import { AuthGuard } from './guards/auth.guard';
import { HomeComponent } from './components/home/home.component';
import { UsersComponent } from './components/users/users.component';
import { UserManageComponent } from './components/user-manage/user-manage.component';
import { CitiesComponent } from './components/cities/cities.component';
import { CityManageComponent } from './components/city-manage/city-manage.component';
import { ProfileComponent } from './components/profile/profile.component';
import { ActivitiesComponent } from './components/activities/activities.component';
import { DirectorsComponent } from './components/directors/directors.component';
import { DirectorManageComponent } from './components/director-manage/director-manage.component';
import { CustomersComponent } from './components/customers/customers.component';
import { CustomerManageComponent } from './components/customer-manage/customer-manage.component';


const routes: Routes = [
	{ path: 'login', component: LoginComponent },
	{ path: '', component: PrivateFrameComponent, canActivate: [AuthGuard],
		children: [
			{ path: '', component: DashboardComponent, canActivateChild: [AuthGuard],
				children:[
					//{ path: '', component: HomeComponent },	
					{ path: '', redirectTo: 'users', pathMatch: 'full' },
					{ path: 'users', component: UsersComponent },	
					{ path: 'users/:id', component: UserManageComponent },	
					{ path: 'home', component: HomeComponent },	
					{ path: 'cities', component: CitiesComponent },	
					{ path: 'cities/:id', component: CityManageComponent },	
					{ path: 'profile', component: ProfileComponent },	
					{ path: 'activities', component: ActivitiesComponent },	
					{ path: 'directors', component: DirectorsComponent },	
					{ path: 'directors/:id', component: DirectorManageComponent },	
					{ path: 'customers', component: CustomersComponent },	
					{ path: 'customers/:id', component: CustomerManageComponent },	
				]
			},
		]
	},
];

@NgModule({
	imports: [RouterModule.forRoot(routes)],
	exports: [RouterModule]
})
export class AppRoutingModule { }
