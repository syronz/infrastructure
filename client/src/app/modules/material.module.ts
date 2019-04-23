import { NgModule } from '@angular/core';
import { MatButtonModule, MatCheckboxModule } from '@angular/material';
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule, MatInputModule } from '@angular/material';
import { MatCardModule } from '@angular/material/card';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatListModule } from '@angular/material/list';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTableModule } from '@angular/material/table';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatSortModule } from '@angular/material/sort';
import { MatSelectModule } from '@angular/material/select';
import { MatSnackBarModule } from '@angular/material/snack-bar';





const modules = [
	MatButtonModule,
	MatCheckboxModule,
	MatIconModule,
	MatFormFieldModule,
	MatInputModule,
	MatCardModule,
	MatToolbarModule,
	MatListModule,
	MatSidenavModule,
	MatProgressSpinnerModule,
	MatTableModule,
	MatPaginatorModule,
	MatSidenavModule,
	MatSortModule,
	MatSelectModule,
	MatSnackBarModule,
];


@NgModule({
  imports: [ ...modules ],
  exports: [ ...modules ],
})
export class MaterialModule { }
