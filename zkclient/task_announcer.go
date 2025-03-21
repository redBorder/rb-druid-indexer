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

func TaskAnnouncer(supervisorTasks []string, tasks []string) []string {
	var tasksToAnnounce []string

	for _, task := range tasks {
		if !taskExists(supervisorTasks, task) {
			tasksToAnnounce = append(tasksToAnnounce, task)
		} else {
		}
	}

	return tasksToAnnounce
}

func taskExists(supervisorTasks []string, task string) bool {
	for _, supervisorTask := range supervisorTasks {
		if supervisorTask == task {
			return true
		}
	}
	return false
}
