/**

Copyright (C) 2012  Roberto Costumero Moreno <roberto@costumero.es>

This file is part of Cosmofs.

Cosmofs is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Cosmofs is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Cosmofs.  If not, see <http://www.gnu.org/licenses/>.

**/

package cosmofs

type chunk struct {
	Name string
	RemPath string
	Owner *Peer
}

type File struct {
	LocalPath string
	GlobalPath string
	Filename string
	Size int64
	Owner *Peer
	Chunks []chunk
	NumChunks int
	Online bool
	KeepCopy bool
	IsDir bool
}

