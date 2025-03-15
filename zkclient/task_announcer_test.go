// Copyright (C) 2025 Eneo Tecnologia S.L.
// Miguel Álvarez <malvarez@redborder.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package zkclient

import (
	"testing"
)

func TestTaskAnnouncer(t *testing.T) {
	tests := []struct {
		name            string
		supervisorTasks []string
		tasks           []string
		expected        []string
	}{
		{
			name:            "No common tasks",
			supervisorTasks: []string{"task1", "task2"},
			tasks:           []string{"task3", "task4"},
			expected:        []string{"task3", "task4"},
		},
		{
			name:            "Some tasks exist",
			supervisorTasks: []string{"task1", "task2"},
			tasks:           []string{"task2", "task3", "task4"},
			expected:        []string{"task3", "task4"},
		},
		{
			name:            "All tasks exist",
			supervisorTasks: []string{"task1", "task2"},
			tasks:           []string{"task1", "task2"},
			expected:        []string{},
		},
		{
			name:            "Empty task lists",
			supervisorTasks: []string{},
			tasks:           []string{"task1", "task2"},
			expected:        []string{"task1", "task2"},
		},
		{
			name:            "Empty supervisor tasks",
			supervisorTasks: []string{},
			tasks:           []string{},
			expected:        []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TaskAnnouncer(tt.supervisorTasks, tt.tasks)
			if !equal(got, tt.expected) {
				t.Errorf("TaskAnnouncer() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
