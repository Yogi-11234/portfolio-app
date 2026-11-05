package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"portfolio/api/internal/domain"
)

var (
	ErrProjectNotFound = errors.New("proyek tidak ditemukan")
	ErrProjectConflict = errors.New("proyek dengan slug ini sudah terdaftar")
)

type ProjectRepository interface {
	List(ctx context.Context, f domain.ProjectFilter) ([]domain.Project, int, error)
	GetByID(ctx context.Context, id string) (*domain.Project, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Project, error)
	Create(ctx context.Context, p *domain.Project, skillIDs []string) error
	Update(ctx context.Context, id string, p *domain.Project, skillIDs []string) error
	UpdateSkills(ctx context.Context, projectID string, skillIDs []string) error
	Delete(ctx context.Context, id string) error
}

type postgresProjectRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresProjectRepository(pool *pgxpool.Pool) ProjectRepository {
	return &postgresProjectRepository{pool: pool}
}

func (r *postgresProjectRepository) List(ctx context.Context, f domain.ProjectFilter) ([]domain.Project, int, error) {
	page := f.Page
	if page < 1 {
		page = 1
	}
	perPage := f.PerPage
	if perPage < 1 || perPage > 100 {
		perPage = 12
	}
	offset := (page - 1) * perPage

	var conditions []string
	var args []interface{}
	argIdx := 1

	if f.Status != "" {
		conditions = append(conditions, fmt.Sprintf("p.status = $%d", argIdx))
		args = append(args, f.Status)
		argIdx++
	}

	if f.FeaturedOnly != nil && *f.FeaturedOnly {
		conditions = append(conditions, fmt.Sprintf("p.is_featured = $%d", argIdx))
		args = append(args, true)
		argIdx++
	}

	if f.CategorySlug != "" {
		conditions = append(conditions, fmt.Sprintf("c.slug = $%d", argIdx))
		args = append(args, f.CategorySlug)
		argIdx++
	}

	if f.SkillSlug != "" {
		conditions = append(conditions, fmt.Sprintf("EXISTS (SELECT 1 FROM project_skills ps JOIN skills s ON ps.skill_id = s.id WHERE ps.project_id = p.id AND s.slug = $%d)", argIdx))
		args = append(args, f.SkillSlug)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT p.id)
		FROM projects p
		LEFT JOIN categories c ON p.category_id = c.id
		%s
	`, whereClause)

	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT p.id, p.category_id, p.title, p.slug, p.summary, p.description_md, p.description_html,
		       p.thumbnail_path, p.repo_url, p.demo_url, p.status, p.is_featured, p.sort_order,
		       p.started_at, p.ended_at, p.created_at, p.updated_at,
		       c.id, c.name, c.slug
		FROM projects p
		LEFT JOIN categories c ON p.category_id = c.id
		%s
		ORDER BY p.is_featured DESC, p.sort_order ASC, p.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIdx, argIdx+1)

	args = append(args, perPage, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var projects []domain.Project
	var projectIDs []string

	for rows.Next() {
		var p domain.Project
		var catID, catName, catSlug *string

		err := rows.Scan(
			&p.ID, &p.CategoryID, &p.Title, &p.Slug, &p.Summary, &p.DescriptionMD, &p.DescriptionHTML,
			&p.ThumbnailPath, &p.RepoURL, &p.DemoURL, &p.Status, &p.IsFeatured, &p.SortOrder,
			&p.StartedAt, &p.EndedAt, &p.CreatedAt, &p.UpdatedAt,
			&catID, &catName, &catSlug,
		)
		if err != nil {
			return nil, 0, err
		}

		if catID != nil && catName != nil && catSlug != nil {
			p.Category = &domain.CategoryBrief{
				ID:   *catID,
				Name: *catName,
				Slug: *catSlug,
			}
		}

		p.Skills = []domain.SkillBrief{}
		projects = append(projects, p)
		projectIDs = append(projectIDs, p.ID)
	}

	if projects == nil {
		projects = []domain.Project{}
	}

	// Fetch skills for these projects
	if len(projectIDs) > 0 {
		skillsMap, err := r.fetchSkillsForProjects(ctx, projectIDs)
		if err == nil {
			for i := range projects {
				if sk, ok := skillsMap[projects[i].ID]; ok {
					projects[i].Skills = sk
				}
			}
		}
	}

	return projects, total, nil
}

func (r *postgresProjectRepository) fetchSkillsForProjects(ctx context.Context, projectIDs []string) (map[string][]domain.SkillBrief, error) {
	query := `
		SELECT ps.project_id, s.id, s.name, s.slug, s.icon_key, s.brand_color
		FROM project_skills ps
		JOIN skills s ON ps.skill_id = s.id
		WHERE ps.project_id = ANY($1)
		ORDER BY s.sort_order ASC, s.name ASC
	`
	rows, err := r.pool.Query(ctx, query, projectIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string][]domain.SkillBrief)
	for rows.Next() {
		var projID string
		var s domain.SkillBrief
		if err := rows.Scan(&projID, &s.ID, &s.Name, &s.Slug, &s.IconKey, &s.BrandColor); err != nil {
			return nil, err
		}
		result[projID] = append(result[projID], s)
	}
	return result, nil
}

func (r *postgresProjectRepository) GetByID(ctx context.Context, id string) (*domain.Project, error) {
	query := `
		SELECT p.id, p.category_id, p.title, p.slug, p.summary, p.description_md, p.description_html,
		       p.thumbnail_path, p.repo_url, p.demo_url, p.status, p.is_featured, p.sort_order,
		       p.started_at, p.ended_at, p.created_at, p.updated_at,
		       c.id, c.name, c.slug
		FROM projects p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1
	`
	var p domain.Project
	var catID, catName, catSlug *string

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.CategoryID, &p.Title, &p.Slug, &p.Summary, &p.DescriptionMD, &p.DescriptionHTML,
		&p.ThumbnailPath, &p.RepoURL, &p.DemoURL, &p.Status, &p.IsFeatured, &p.SortOrder,
		&p.StartedAt, &p.EndedAt, &p.CreatedAt, &p.UpdatedAt,
		&catID, &catName, &catSlug,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	if catID != nil && catName != nil && catSlug != nil {
		p.Category = &domain.CategoryBrief{
			ID:   *catID,
			Name: *catName,
			Slug: *catSlug,
		}
	}

	skillsMap, _ := r.fetchSkillsForProjects(ctx, []string{p.ID})
	if sk, ok := skillsMap[p.ID]; ok {
		p.Skills = sk
	} else {
		p.Skills = []domain.SkillBrief{}
	}

	return &p, nil
}

func (r *postgresProjectRepository) GetBySlug(ctx context.Context, slug string) (*domain.Project, error) {
	query := `
		SELECT p.id, p.category_id, p.title, p.slug, p.summary, p.description_md, p.description_html,
		       p.thumbnail_path, p.repo_url, p.demo_url, p.status, p.is_featured, p.sort_order,
		       p.started_at, p.ended_at, p.created_at, p.updated_at,
		       c.id, c.name, c.slug
		FROM projects p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.slug = $1
	`
	var p domain.Project
	var catID, catName, catSlug *string

	err := r.pool.QueryRow(ctx, query, slug).Scan(
		&p.ID, &p.CategoryID, &p.Title, &p.Slug, &p.Summary, &p.DescriptionMD, &p.DescriptionHTML,
		&p.ThumbnailPath, &p.RepoURL, &p.DemoURL, &p.Status, &p.IsFeatured, &p.SortOrder,
		&p.StartedAt, &p.EndedAt, &p.CreatedAt, &p.UpdatedAt,
		&catID, &catName, &catSlug,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	if catID != nil && catName != nil && catSlug != nil {
		p.Category = &domain.CategoryBrief{
			ID:   *catID,
			Name: *catName,
			Slug: *catSlug,
		}
	}

	skillsMap, _ := r.fetchSkillsForProjects(ctx, []string{p.ID})
	if sk, ok := skillsMap[p.ID]; ok {
		p.Skills = sk
	} else {
		p.Skills = []domain.SkillBrief{}
	}

	return &p, nil
}

func (r *postgresProjectRepository) Create(ctx context.Context, p *domain.Project, skillIDs []string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO projects (category_id, title, slug, summary, description_md, description_html,
		                      thumbnail_path, repo_url, demo_url, status, is_featured, sort_order,
		                      started_at, ended_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, created_at, updated_at
	`
	err = tx.QueryRow(ctx, query,
		p.CategoryID, p.Title, p.Slug, p.Summary, p.DescriptionMD, p.DescriptionHTML,
		p.ThumbnailPath, p.RepoURL, p.DemoURL, p.Status, p.IsFeatured, p.SortOrder,
		p.StartedAt, p.EndedAt,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return err
	}

	// Insert project_skills
	if len(skillIDs) > 0 {
		for _, sID := range skillIDs {
			_, err := tx.Exec(ctx, `INSERT INTO project_skills (project_id, skill_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, p.ID, sID)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

func (r *postgresProjectRepository) Update(ctx context.Context, id string, p *domain.Project, skillIDs []string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		UPDATE projects
		SET category_id = $1, title = $2, slug = $3, summary = $4, description_md = $5,
		    description_html = $6, thumbnail_path = $7, repo_url = $8, demo_url = $9,
		    status = $10, is_featured = $11, sort_order = $12, started_at = $13,
		    ended_at = $14, updated_at = now()
		WHERE id = $15
		RETURNING updated_at
	`
	err = tx.QueryRow(ctx, query,
		p.CategoryID, p.Title, p.Slug, p.Summary, p.DescriptionMD,
		p.DescriptionHTML, p.ThumbnailPath, p.RepoURL, p.DemoURL,
		p.Status, p.IsFeatured, p.SortOrder, p.StartedAt,
		p.EndedAt, id,
	).Scan(&p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrProjectNotFound
		}
		return err
	}

	// Replace project_skills
	_, err = tx.Exec(ctx, `DELETE FROM project_skills WHERE project_id = $1`, id)
	if err != nil {
		return err
	}

	if len(skillIDs) > 0 {
		for _, sID := range skillIDs {
			_, err := tx.Exec(ctx, `INSERT INTO project_skills (project_id, skill_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, id, sID)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

func (r *postgresProjectRepository) UpdateSkills(ctx context.Context, projectID string, skillIDs []string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Ensure project exists
	var exists bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM projects WHERE id = $1)`, projectID).Scan(&exists); err != nil || !exists {
		return ErrProjectNotFound
	}

	_, err = tx.Exec(ctx, `DELETE FROM project_skills WHERE project_id = $1`, projectID)
	if err != nil {
		return err
	}

	for _, sID := range skillIDs {
		_, err := tx.Exec(ctx, `INSERT INTO project_skills (project_id, skill_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, projectID, sID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *postgresProjectRepository) Delete(ctx context.Context, id string) error {
	cmdTag, err := r.pool.Exec(ctx, `DELETE FROM projects WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrProjectNotFound
	}
	return nil
}
