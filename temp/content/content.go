// content created for reading the modules data from the disk and write into MongoDB
// modules_dir
// ∟ stage-<id>
//    ∟ module-<id>
//    |  ∟ module.json
//    |  ∟ part-<id>
//    |  |  ∟ data.html
//    |  ∟ part-<id>
//    |	 ∟ data.html
//    |	 ∟ questions.json
//    ∟ module-<id>
//    ...
package content

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/CyrilKuzmin/itpath69/internal/domain/module"
	"github.com/CyrilKuzmin/itpath69/internal/domain/tests"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ContentManager struct {
	moduleService module.Service
	testService   tests.Service
	log           *zap.Logger
	ModulesTotal  int
}

func NewContentManager(ms module.Service, t tests.Service, log *zap.Logger) *ContentManager {
	return &ContentManager{
		moduleService: ms,
		testService:   t,
		log:           log,
	}
}

func (cm *ContentManager) UpdateContentinStorage() {
	var modules []module.Module
	var qs []tests.Question
	baseDir := "static/modules"
	baseDirInfo, err := os.ReadDir(baseDir)
	if err != nil {
		log.Fatal("cannot list base dir", zap.Error(err), zap.String("dir", baseDir))
	}
	for _, m := range baseDirInfo {
		modulePath := fmt.Sprintf("%v/%v", baseDir, m.Name())
		m := cm.readModule(modulePath)
		moduleQs := cm.readQuestions(modulePath)
		qs = append(qs, moduleQs...)
		m.Data = make([]module.Part, 0)
		moduleDirInfo, err := os.ReadDir(modulePath)
		parts := make([]module.Part, 0)
		if err != nil {
			log.Fatal("cannot list module dir", zap.Error(err), zap.String("dir", modulePath))
		}
		for _, p := range moduleDirInfo {
			if !p.IsDir() {
				continue
			}
			partPath := fmt.Sprintf("%v/%v", modulePath, p.Name())
			parts = append(parts, cm.readPart(partPath))
		}
		m.Data = append(m.Data, parts...)
		modules = append(modules, m)
	}
	err = cm.moduleService.CreateModules(context.Background(), modules)
	if err != nil {
		log.Fatal("cannot insert modules", zap.Error(err))
	}
	// I don't want to re-generate new uuids for questions during the development
	// err = cm.testService.SaveQuestions(context.Background(), qs)
	// if err != nil {
	// 	log.Fatal("cannot insert questions", zap.Error(err))
	// }
	cm.ModulesTotal = len(modules)
}

func (cm *ContentManager) readModule(moduleDir string) module.Module {
	var m module.Module
	var moduleMeta module.ModuleMeta
	meta, err := os.Open(fmt.Sprintf("%v/module.json", moduleDir))
	defer meta.Close()
	if err != nil {
		log.Fatal("cannot read a module", zap.Error(err), zap.String("module", moduleDir))
	}
	metaValue, _ := io.ReadAll(meta)
	json.Unmarshal([]byte(metaValue), &moduleMeta)
	m.Meta = moduleMeta
	m.Id = moduleMeta.Id
	return m
}

func (cm *ContentManager) readPart(partDir string) module.Part {
	var part module.Part
	meta, err := os.Open(fmt.Sprintf("%v/part.json", partDir))
	defer meta.Close()
	if err != nil {
		log.Fatal("cannot read a part", zap.Error(err), zap.String("part", partDir))
	}
	metaValue, _ := io.ReadAll(meta)
	json.Unmarshal([]byte(metaValue), &part)
	data, err := os.Open(fmt.Sprintf("%v/data.html", partDir))
	defer data.Close()
	if err != nil {
		log.Fatal("cannot read a part data", zap.Error(err), zap.String("part", partDir))
	}
	dataValue, _ := io.ReadAll(data)
	part.Data = string(dataValue)
	return part
}

func (cm *ContentManager) readQuestions(moduleDir string) []tests.Question {
	questions := make([]tests.Question, 0)
	filename := fmt.Sprintf("%v/questions.json", moduleDir)
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		log.Fatal("cannot read test questions", zap.Error(err), zap.String("file", filename))
	}
	data, _ := io.ReadAll(f)
	json.Unmarshal([]byte(data), &questions)
	for i, _ := range questions {
		questions[i].Id = uuid.NewString()
	}
	return questions
}
