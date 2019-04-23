import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';
import { HTTP_INTERCEPTORS } from '@angular/common/http';
import { JwtInterceptor } from './interceptors/jwt.interceptor';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { LoginComponent } from './components/login/login.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MaterialModule } from './modules/material.module';
import { DictionaryPipe } from './pipes/dictionary.pipe';
import { DashboardComponent } from './components/dashboard/dashboard.component';
import { PrivateFrameComponent } from './components/private-frame/private-frame.component';
import { NavbarComponent } from './components/navbar/navbar.component';
import { SidebarComponent } from './components/sidebar/sidebar.component';
import { HomeComponent } from './components/home/home.component';
import { UsersComponent } from './components/users/users.component';
import { UserManageComponent } from './components/user-manage/user-manage.component';
import { LoadingSpinnerComponent } from './components/loading-spinner/loading-spinner.component';
import { CitiesComponent } from './components/cities/cities.component';
import { CityManageComponent } from './components/city-manage/city-manage.component';
import { ProfileComponent } from './components/profile/profile.component';
import { ActivitiesComponent } from './components/activities/activities.component';
import { DirectorsComponent } from './components/directors/directors.component';
import { DirectorManageComponent } from './components/director-manage/director-manage.component';

@NgModule({
  declarations: [
    AppComponent,
    LoginComponent,
    DictionaryPipe,
    DashboardComponent,
    PrivateFrameComponent,
    NavbarComponent,
    SidebarComponent,
    HomeComponent,
    UsersComponent,
    UserManageComponent,
		LoadingSpinnerComponent,
		CitiesComponent,
		CityManageComponent,
		ProfileComponent,
		ActivitiesComponent,
		DirectorsComponent,
		DirectorManageComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    BrowserAnimationsModule,
		MaterialModule,
		FormsModule,
		ReactiveFormsModule,
		HttpClientModule,
	],
	entryComponents: [
	],


	providers: [
		{ provide: HTTP_INTERCEPTORS, useClass: JwtInterceptor, multi: true },
	],
	bootstrap: [AppComponent]
})
export class AppModule { }
