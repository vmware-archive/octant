import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { Filter, LabelFilterService } from 'src/app/services/label-filter.service';

@Component({
  selector: 'app-filters',
  templateUrl: './filters.component.html',
  styleUrls: ['./filters.component.scss'],
})
export class FiltersComponent implements OnInit {
  filters: Filter[];

  constructor(
    private labelFilter: LabelFilterService,
    private router: Router,
    private activatedRoute: ActivatedRoute
  ) {}

  ngOnInit() {
    this.labelFilter.filters().subscribe((filters) => {
      this.filters = filters;
      const filterParams = filters.map((filter) => encodeURIComponent(`${filter.key}:${filter.value}`));
      const queryParams: Params = {
        filter: filterParams,
      };

      this.router.navigate([], {
        relativeTo: this.activatedRoute,
        queryParams,
        queryParamsHandling: 'merge',
      });
    });
  }

  remove(filter: Filter) {
    this.labelFilter.remove(filter);
  }
}
