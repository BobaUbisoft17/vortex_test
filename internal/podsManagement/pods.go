package podsmanagement

import (
	"errors"
	"slices"
)

type PodList struct {
	pods []string
}

func (p *PodList) CreatePod(name string) error {
	if slices.Contains[[]string, string](p.pods, name) {
		return errors.New("pod is already created")
	}

	p.pods = append(p.pods, name)
	return nil
}

func (p *PodList) DeletePod(name string) error {
	if !slices.Contains[[]string, string](p.pods, name) {
		return errors.New("pod does not exist")
	}

	for i := range p.pods {
		if p.pods[i] == name {
			p.pods = append(p.pods[:i], p.pods[i+1:]...)
			break
		}
	}
	return nil
}

func (p *PodList) GetPodList() ([]string, error) {
	return p.pods, nil
}
