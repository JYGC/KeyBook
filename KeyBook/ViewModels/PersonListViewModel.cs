﻿using KeyBook.Models;

namespace KeyBook.ViewModels
{
    public class PersonListViewModel
    {
        public List<Person> Persons { get; set; }
        public Dictionary<int, string> PersonTypes { get; set; }
    }
}
